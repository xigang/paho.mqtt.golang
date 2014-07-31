/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

// Package mqtt provides an MQTT v3.1 client library.
package mqtt

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"time"
)

// MqttClient is a lightweight MQTT v3.1 Client for communicating
// with an MQTT server using non-blocking methods that allow work
// to be done in the background.

// An application may connect to an MQTT server using:
//   A plain TCP socket
//   A secure SSL/TLS socket
//   A websocket

// To enable ensured message delivery at Quality of Service (QoS) levels
// described in the MQTT spec, a message persistence mechanism must be
// used. This is done by providing a type which implements the Store
// interface. For convenience, FileStore and MemoryStore are provided
// implementations that should be sufficient for most use cases. More
// information can be found in their respective documentation.

// Numerous connection options may be specified by configuring a
// and then supplying a ClientOptions type.
type MqttClient struct {
	sync.RWMutex
	conn            net.Conn
	bufferedConn    *bufio.ReadWriter
	ibound          chan *Message
	obound          chan sendable
	oboundP         chan *Message
	begin           chan ConnRC
	errors          chan error
	stop            chan struct{}
	receipts        *receiptMap
	t               *Tracer
	persist         Store
	options         ClientOptions
	lastContact     lastcontact
	pingOutstanding bool
	connected       bool
}

// NewClient will create an MQTT v3.1 client with all of the options specified
// in the provided ClientOptions. The client must have the Start method called
// on it before it may be used. This is to make sure resources (such as a net
// connection) are created before the application is actually ready.
func NewClient(ops *ClientOptions) *MqttClient {
	c := &MqttClient{}
	c.options = *ops

	if c.options.store == nil {
		c.options.store = NewMemoryStore()
	}
	c.persist = c.options.store
	c.connected = false
	return c
}

func (c *MqttClient) IsConnected() bool {
	c.RLock()
	defer c.RUnlock()
	return c.connected
}

// Start will create a connection to the message broker
// If clean session is false, then a slice will
// be returned containing Receipts for all messages
// that were in-flight at the last disconnect.
// If clean session is true, then any existing client
// state will be removed.
func (c *MqttClient) Start() ([]Receipt, error) {

	c.t = &Tracer{
		c.options.tracelevel,
		c.options.tracefile,
		c.options.clientId,
	}

	c.options.store.SetTracer(c.t)

	c.trace_v(CLI, "Start()")

	for _, broker := range c.options.servers {
		conn, err := openConnection(broker, c.options.tlsconfig)
		if err == nil {
			c.conn = conn
			c.trace_v(CLI, "connected to broker")
			break
		} else {
			c.trace_w(CLI, "failed to connect to broker, trying next")
		}
	}

	if c.conn == nil {
		c.trace_e(CLI, "Failed to connect to a broker")
		return nil, errors.New("Failed to connect to a broker")
	}
	c.bufferedConn = bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

	c.persist.Open()
	c.receipts = newReceiptMap()

	c.trace_v(CLI, "about to start generateMsgIds")
	c.options.mids.generateMsgIds()

	c.obound = make(chan sendable)
	c.ibound = make(chan *Message)
	c.oboundP = make(chan *Message)
	c.errors = make(chan error)
	c.stop = make(chan struct{})

	go outgoing(c)
	go alllogic(c)

	cm := newConnectMsgFromOptions(c.options)
	c.trace_v(CLI, "about to write new connect msg")
	c.oboundP <- cm

	rc := connect(c)
	if rc != CONN_ACCEPTED {
		c.trace_c(CLI, "CONNACK was not CONN_ACCEPTED, but rather %s", rc2str(rc))
		// Stop all go routines except outgoing
		close(c.stop)
		c.conn.Close()
		return nil, chkrc(rc)
	}

	c.options.pubChanZero = make(chan *Message, 1000)
	c.options.pubChanOne = make(chan *Message, 1000)
	c.options.pubChanTwo = make(chan *Message, 1000)
	c.options.msgRouter.matchAndDispatch(c.options.pubChanZero, c.options.order, c)
	c.options.msgRouter.matchAndDispatch(c.options.pubChanOne, c.options.order, c)
	c.options.msgRouter.matchAndDispatch(c.options.pubChanTwo, c.options.order, c)

	c.connected = true
	c.trace_v(CLI, "client is connected")

	if c.options.timeout != 0 {
		go keepalive(c)
	}

	// Take care of any messages in the store
	var leftovers []Receipt
	if c.options.cleanses == false {
		leftovers = c.resume()
	} else {
		c.persist.Reset()
	}

	// Do not start incoming until resume has completed
	go incoming(c)

	c.trace_v(CLI, "exit startMqttClient")
	if chkrc(rc) != nil {
		// Cleanup before returning.
		close(c.stop)
		c.conn.Close()
	}
	return leftovers, chkrc(rc)
}

// Disconnect will end the connection with the server, but not before waiting
// the specified number of milliseconds to wait for existing work to be
// completed.
func (c *MqttClient) Disconnect(quiesce uint) {
	if !c.IsConnected() {
		c.trace_w(CLI, "already disconnected")
		return
	}
	c.trace_v(CLI, "disconnecting")
	c.connected = false

	// wait for work to finish, or quiesce time consumed
	end := time.After(time.Duration(quiesce) * time.Millisecond)

	// for now we just wait for the time specified and hope the work is done
	select {
	case <-end:
		c.trace_v(CLI, "quiesce expired, forcing disconnect")
		// case <- other:
		// 	c.trace_v(CLI, "finished processing work, graceful disconnect")
	}
	c.disconnect()
}

// ForceDisconnect will end the connection with the mqtt broker immediately.
func (c *MqttClient) ForceDisconnect() {
	if !c.IsConnected() {
		c.trace_w(CLI, "already disconnected")
		return
	}
	c.trace_v(CLI, "forcefully disconnecting")
	c.disconnect()
}

func (c *MqttClient) disconnect() {
	c.connected = false
	dm := newDisconnectMsg()

	// Stop all go routines except outgoing
	close(c.stop)

	// Send disconnect message and stop outgoing
	c.oboundP <- dm

	c.trace_v(CLI, "disconnected")
	c.persist.Close()
}

// Publish will publish a message with the specified QoS
// and content to the specified topic.
// Returns a read only channel used to track
// the delivery of the message.
func (c *MqttClient) Publish(qos QoS, topic string, payload interface{}) <-chan Receipt {
	var pub *Message
	switch payload.(type) {
	case string:
		pub = newPublishMsg(qos, topic, []byte(payload.(string)))
	case []byte:
		pub = newPublishMsg(qos, topic, payload.([]byte))
	default:
		return nil
	}

	r := make(chan Receipt, 1)
	c.trace_v(CLI, "sending publish message, topic: %s", topic)

	select {
	case c.obound <- sendable{pub, r}:
	case <-time.After(time.Second):
		close(r)
	}
	return r
}

// PublishMessage will publish a Message to the specified topic.
// Returns a read only channel used to track
// the delivery of the message.
func (c *MqttClient) PublishMessage(topic string, message *Message) <-chan Receipt {
	// Just reuse pieces from the existing message
	// so that message id etc aren't set
	pub := newPublishMsg(message.QoS(), topic, message.payload)
	pub.SetRetainedFlag(message.RetainedFlag())

	r := make(chan Receipt, 1)

	c.trace_v(CLI, "sending publish message, topic: %s", topic)

	select {
	case c.obound <- sendable{pub, r}:
		return r
	case <-time.After(time.Second):
		close(r)
		return nil
	}
}

// Start a new subscription. Provide a MessageHandler to be executed when
// a message is published on one of the topics provided.
func (c *MqttClient) StartSubscription(callback MessageHandler, filters ...*TopicFilter) (<-chan Receipt, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}
	c.trace_v(CLI, "enter StartSubscription")
	submsg := newSubscribeMsg(filters...)
	chkcond(submsg != nil)

	if callback != nil {
		for i := range filters {
			c.options.msgRouter.addRoute(filters[i].string, callback)
		}
	}

	r := make(chan Receipt, 1)

	c.obound <- sendable{submsg, r}

	c.trace_v(CLI, "exit StartSubscription")
	return r, nil
}

// EndSubscription will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *MqttClient) EndSubscription(topics ...string) (<-chan Receipt, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}
	c.trace_v(CLI, "enter EndSubscription")
	usmsg := newUnsubscribeMsg(topics...)

	r := make(chan Receipt, 1)

	c.obound <- sendable{usmsg, r}
	for _, topic := range topics {
		c.options.msgRouter.deleteRoute(topic)
	}

	c.trace_v(CLI, "exit EndSubscription")
	return r, nil
}
