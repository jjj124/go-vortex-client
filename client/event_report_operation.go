package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
	"maps"
)

type EventReportReply struct {
	msg.ReceivedMsg
}

type EventReportOperation interface {
	NeedReply() EventReportOperation
	WithValue(key string, val any) EventReportOperation
	WithTs(ts int64) EventReportOperation
	Execute() futures.Future[EventReportReply]
}

type defaultEventReportOperation struct {
	adapter   AdapterClient
	device    futures.Future[Device]
	v         map[string]any
	needReply bool
	ts        int64
}

func (d *defaultEventReportOperation) WithTs(ts int64) EventReportOperation {
	d.ts = ts
	return d
}

func NewEventReportOperation(adapter AdapterClient, device futures.Future[Device]) EventReportOperation {
	return &defaultEventReportOperation{
		adapter,
		device,
		make(map[string]any, 0),
		false,
		-1,
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
		var baseMsg = msg.NewBaseMsg().WithMethod(DeviceEventReport)
		if d.ts > -1 {
			baseMsg.WithTs(d.ts)
		}
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
