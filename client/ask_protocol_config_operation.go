package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
)

type AskConfigReply struct {
	*msg.ReceivedMsg
}

func (a *AskConfigReply) GetContent() string {
	var content = a.Payload()["content"]
	if content != nil {
		return content.(string)
	} else {
		return ""
	}
}
func (a *AskConfigReply) GetFormat() string {
	var format = a.Payload()["format"]
	if format != nil {
		return format.(string)
	} else {
		return ""
	}
}

type AskProtocolConfigOperation interface {
	Execute() futures.Future[AskConfigReply]
}

type askProtocolConfigOperation struct {
	adapterClient AdapterClient
}

func (a *askProtocolConfigOperation) Execute() futures.Future[AskConfigReply] {
	var baseMsg = msg.NewBaseMsg()
	baseMsg.WithMethod("vortex.adapter.ask-protocol-config")
	var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
	var f = a.adapterClient.Delivery(deliveryMsg)
	return futures.Then(f, func(m *msg.ReceivedMsg) *AskConfigReply {
		return &AskConfigReply{
			m,
		}
	})
}
func NewAskConfigOperation(adapter AdapterClient) AskProtocolConfigOperation {
	return &askProtocolConfigOperation{
		adapterClient: adapter,
	}
}
