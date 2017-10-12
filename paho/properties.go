package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

type Properties packets.IDValuePair

func NewProperties() *Properties {
	return &Properties{}
}

func (p *Properties) SetPayloadFormat(v byte) *Properties {
	p.PayloadFormat = &v
	return p
}

func (p *Properties) SetPubExpiry(v uint32) *Properties {
	p.PubExpiry = &v
	return p
}

func (p *Properties) SetContentType(v string) *Properties {
	p.ContentType = v
	return p
}

func (p *Properties) SetReplyTopic(v string) *Properties {
	p.ReplyTopic = v
	return p
}

func (p *Properties) SetCorrelationData(v []byte) *Properties {
	p.CorrelationData = v
	return p
}

func (p *Properties) SetSubscriptionIdentifier(v *uint32) *Properties {
	p.SubscriptionIdentifier = v
	return p
}

func (p *Properties) SetSessionExpiryInterval(v uint32) *Properties {
	p.SessionExpiryInterval = &v
	return p
}

func (p *Properties) SetAssignedClientID(v string) *Properties {
	p.AssignedClientID = v
	return p
}

func (p *Properties) SetServerKeepAlive(v uint16) *Properties {
	p.ServerKeepAlive = &v
	return p
}

func (p *Properties) SetAuthMethod(v string) *Properties {
	p.AuthMethod = v
	return p
}

func (p *Properties) SetAuthData(v []byte) *Properties {
	p.AuthData = v
	return p
}

func (p *Properties) SetRequestProblemInfo(v byte) *Properties {
	p.RequestProblemInfo = &v
	return p
}

func (p *Properties) SetWillDelayInterval(v uint32) *Properties {
	p.WillDelayInterval = &v
	return p
}

func (p *Properties) SetRequestResponseInfo(v byte) *Properties {
	p.RequestResponseInfo = &v
	return p
}

func (p *Properties) SetResponseInfo(v string) *Properties {
	p.ResponseInfo = v
	return p
}

func (p *Properties) SetServerReference(v string) *Properties {
	p.ServerReference = v
	return p
}

func (p *Properties) SetReasonString(v string) *Properties {
	p.ReasonString = v
	return p
}

func (p *Properties) SetReceiveMaximum(v uint16) *Properties {
	p.ReceiveMaximum = &v
	return p
}

func (p *Properties) SetTopicAliasMaximum(v uint16) *Properties {
	p.TopicAliasMaximum = &v
	return p
}

func (p *Properties) SetTopicAlias(v uint16) *Properties {
	p.TopicAlias = &v
	return p
}

func (p *Properties) SetMaximumQOS(v byte) *Properties {
	p.MaximumQOS = &v
	return p
}

func (p *Properties) SetRetainAvailable(v byte) *Properties {
	p.RetainAvailable = &v
	return p
}

func (p *Properties) SetUserProperty(k, v string) *Properties {
	p.UserProperty[k] = v
	return p
}

func (p *Properties) SetMaximumPacketSize(v uint32) *Properties {
	p.MaximumPacketSize = &v
	return p
}

func (p *Properties) SetWildcardSubAvailable(v byte) *Properties {
	p.WildcardSubAvailable = &v
	return p
}

func (p *Properties) SetSubIDAvailable(v byte) *Properties {
	p.SubIDAvailable = &v
	return p
}

func (p *Properties) SetSharedSubAvailable(v byte) *Properties {
	p.SharedSubAvailable = &v
	return p
}

func (p *Properties) Validate(pt packets.PacketType) (bool, []string) {
	if p == nil {
		return true, nil
	}
	valid := true
	var invalid []string
	if p.PayloadFormat != nil && !packets.ValidateID(pt, packets.IDVPPayloadFormat) {
		invalid = append(invalid, "PayloadFormat")
		valid = false
	}
	if p.PubExpiry != nil && !packets.ValidateID(pt, packets.IDVPPubExpiry) {
		invalid = append(invalid, "PubExpiry")
		valid = false
	}
	if p.ContentType != "" && !packets.ValidateID(pt, packets.IDVPContentType) {
		invalid = append(invalid, "ContentType")
		valid = false
	}
	if p.ReplyTopic != "" && !packets.ValidateID(pt, packets.IDVPReplyTopic) {
		invalid = append(invalid, "ReplyTopic")
		valid = false
	}
	if p.CorrelationData != nil && !packets.ValidateID(pt, packets.IDVPCorrelationData) {
		invalid = append(invalid, "CorrelationData")
		valid = false
	}
	if p.SubscriptionIdentifier != nil && !packets.ValidateID(pt, packets.IDVPSubscriptionIdentifier) {
		invalid = append(invalid, "SubscriptionIdentifier")
		valid = false
	}
	if p.SessionExpiryInterval != nil && !packets.ValidateID(pt, packets.IDVPSessionExpiryInterval) {
		invalid = append(invalid, "SessionExpiryInterval")
		valid = false
	}
	if p.AssignedClientID != "" && !packets.ValidateID(pt, packets.IDVPAssignedClientID) {
		invalid = append(invalid, "AssignedClientID")
		valid = false
	}
	if p.ServerKeepAlive != nil && !packets.ValidateID(pt, packets.IDVPServerKeepAlive) {
		invalid = append(invalid, "ServerKeepAlive")
		valid = false
	}
	if p.AuthMethod != "" && !packets.ValidateID(pt, packets.IDVPAuthMethod) {
		invalid = append(invalid, "AuthMethod")
		valid = false
	}
	if p.AuthData != nil && !packets.ValidateID(pt, packets.IDVPAuthData) {
		invalid = append(invalid, "AuthData")
		valid = false
	}
	if p.RequestProblemInfo != nil && !packets.ValidateID(pt, packets.IDVPRequestProblemInfo) {
		invalid = append(invalid, "RequestProblemInfo")
		valid = false
	}
	if p.WillDelayInterval != nil && !packets.ValidateID(pt, packets.IDVPWillDelayInterval) {
		invalid = append(invalid, "WillDelayInterval")
		valid = false
	}
	if p.RequestResponseInfo != nil && !packets.ValidateID(pt, packets.IDVPRequestResponseInfo) {
		invalid = append(invalid, "RequestResponseInfo")
		valid = false
	}
	if p.ResponseInfo != "" && !packets.ValidateID(pt, packets.IDVPResponseInfo) {
		invalid = append(invalid, "ResponseInfo")
		valid = false
	}
	if p.ServerReference != "" && !packets.ValidateID(pt, packets.IDVPServerReference) {
		invalid = append(invalid, "ServerReference")
		valid = false
	}
	if p.ReasonString != "" && !packets.ValidateID(pt, packets.IDVPReasonString) {
		invalid = append(invalid, "ReasonString")
		valid = false
	}
	if p.ReceiveMaximum != nil && !packets.ValidateID(pt, packets.IDVPReceiveMaximum) {
		invalid = append(invalid, "ReceiveMaximum")
		valid = false
	}
	if p.TopicAliasMaximum != nil && !packets.ValidateID(pt, packets.IDVPTopicAliasMaximum) {
		invalid = append(invalid, "TopicAliasMaximum")
		valid = false
	}
	if p.TopicAlias != nil && !packets.ValidateID(pt, packets.IDVPTopicAlias) {
		invalid = append(invalid, "TopicAlias")
		valid = false
	}
	if p.MaximumQOS != nil && !packets.ValidateID(pt, packets.IDVPMaximumQOS) {
		invalid = append(invalid, "MaximumQOS")
		valid = false
	}
	if p.RetainAvailable != nil && !packets.ValidateID(pt, packets.IDVPRetainAvailable) {
		invalid = append(invalid, "RetainAvailable")
		valid = false
	}
	if p.UserProperty != nil && !packets.ValidateID(pt, packets.IDVPUserProperty) {
		invalid = append(invalid, "UserProperty")
		valid = false
	}
	if p.MaximumPacketSize != nil && !packets.ValidateID(pt, packets.IDVPMaximumPacketSize) {
		invalid = append(invalid, "MaximumPacketSize")
		valid = false
	}
	if p.WildcardSubAvailable != nil && !packets.ValidateID(pt, packets.IDVPWildcardSubAvailable) {
		invalid = append(invalid, "WildcardSubAvailable")
		valid = false
	}
	if p.SubIDAvailable != nil && !packets.ValidateID(pt, packets.IDVPSubIDAvailable) {
		invalid = append(invalid, "SubIDAvailable")
		valid = false
	}
	if p.SharedSubAvailable != nil && !packets.ValidateID(pt, packets.IDVPSharedSubAvailable) {
		invalid = append(invalid, "SharedSubAvailable")
		valid = false
	}

	return valid, invalid
}
