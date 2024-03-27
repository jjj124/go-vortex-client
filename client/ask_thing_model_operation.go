package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
)

type AskThingModelReply struct {
	*msg.ReceivedMsg
}

func (a *AskThingModelReply) GetProps() []ThingModel {
	var payload = a.Payload()
	var ret = make([]ThingModel, 0)
	for key, val := range payload {
		var item = val.(map[string]any)
		var modelType = item["type"].(string)
		var name = item["name"].(string)
		var dataType = item["data_type"].(string)
		var thingModel = ThingModel{
			Name:     name,
			DataType: dataType,
			Symbol:   key,
		}
		if modelType == "prop" {
			ret = append(ret, thingModel)
		}
	}
	return ret
}

type AskThingModelOperation interface {
	WithPid(pid string) AskThingModelOperation
	Execute() futures.Future[AskThingModelReply]
}

type askThingModelOperation struct {
	adapterClient AdapterClient
	pid           string
}

func (a *askThingModelOperation) WithPid(pid string) AskThingModelOperation {
	a.pid = pid
	return a
}

func (a *askThingModelOperation) Execute() futures.Future[AskThingModelReply] {
	var baseMsg = msg.NewBaseMsg()
	baseMsg.WithMethod("vortex.adapter.ask-thing-model-definition")
	var payload = baseMsg.Payload()
	if a.pid != "" {
		payload["pid"] = a.pid
	}
	var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
	var f = a.adapterClient.Delivery(deliveryMsg)
	return futures.Then(f, func(m *msg.ReceivedMsg) *AskThingModelReply {
		var v = AskThingModelReply{m}
		return &v
	})
}

func NewAskThingModelOperation(adapter AdapterClient) AskThingModelOperation {
	return &askThingModelOperation{
		adapterClient: adapter,
	}
}
