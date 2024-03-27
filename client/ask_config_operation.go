package client

import (
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	futures "github.com/jjj124/go-future"
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

type AskConfigOperation interface {
	Execute() futures.Future[AskConfigReply]
}

type askConfigOperation struct {
	adapterClient AdapterClient
}

func (a *askConfigOperation) Execute() futures.Future[AskConfigReply] {
	var baseMsg = msg.NewBaseMsg()
	baseMsg.WithMethod("vortex.adapter.ask-config")
	var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
	var f = a.adapterClient.Delivery(deliveryMsg)
	return futures.Then(f, func(m *msg.ReceivedMsg) *AskConfigReply {
		return &AskConfigReply{
			m,
		}
	})
}
func NewAskConfigOperation(adapter AdapterClient) AskConfigOperation {
	return &askConfigOperation{
		adapterClient: adapter,
	}
}
