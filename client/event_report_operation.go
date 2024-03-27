package client

import (
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	futures "github.com/jjj124/go-future"
	"maps"
)

type EventReportReply struct {
	msg.ReceivedMsg
}

type EventReportOperation interface {
	NeedReply() EventReportOperation
	WithValue(key string, val any) EventReportOperation
	Execute() futures.Future[EventReportReply]
}

type defaultEventReportOperation struct {
	adapter   AdapterClient
	device    futures.Future[Device]
	v         map[string]any
	needReply bool
}

func NewEventReportOperation(adapter AdapterClient, device futures.Future[Device]) EventReportOperation {
	return &defaultEventReportOperation{
		adapter,
		device,
		make(map[string]any, 0),
		false,
	}
}

func (d *defaultEventReportOperation) NeedReply() EventReportOperation {
	d.needReply = true
	return d
}

func (d *defaultEventReportOperation) WithValue(key string, val any) EventReportOperation {
	d.v[key] = val
	return d
}

func (d *defaultEventReportOperation) Execute() futures.Future[EventReportReply] {
	var ret = futures.NewFuture[EventReportReply]()
	d.device.WhenComplete(func(dev *Device, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			return
		}
		var baseMsg = msg.NewBaseMsg().WithMethod("device.event.report")
		var payload = baseMsg.Payload()
		maps.Copy(payload, d.v)
		if d.needReply {
			baseMsg.Header()["need_reply"] = true
		}
		baseMsg.Header()["adapter_pid"] = (*dev).Pid()
		baseMsg.Header()["adapter_did"] = (*dev).Did()
		var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
		d.adapter.Delivery(deliveryMsg).WhenComplete(func(m *msg.ReceivedMsg, err error) {
			if err != nil {
				ret.CompleteExceptionally(err)
			} else {
				var v = EventReportReply{*m}
				ret.Complete(&v)
			}
		})

	})
	return ret
}
