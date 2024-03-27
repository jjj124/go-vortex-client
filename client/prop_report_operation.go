package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
	"maps"
)

type PropReportReply struct {
	msg.ReceivedMsg
}

type PropReportOperation interface {
	NeedReply() PropReportOperation
	WithValue(key string, val any) PropReportOperation
	Execute() futures.Future[PropReportReply]
}

type defaultPropReportOperation struct {
	adapter   AdapterClient
	device    futures.Future[Device]
	v         map[string]any
	needReply bool
}

func NewPropReportOperation(adapter AdapterClient, device futures.Future[Device]) PropReportOperation {
	return &defaultPropReportOperation{
		adapter,
		device,
		make(map[string]any, 0),
		false,
	}
}

func (d *defaultPropReportOperation) NeedReply() PropReportOperation {
	d.needReply = true
	return d
}

func (d *defaultPropReportOperation) WithValue(key string, val any) PropReportOperation {
	d.v[key] = val
	return d
}

func (d *defaultPropReportOperation) Execute() futures.Future[PropReportReply] {

	var ret = futures.NewFuture[PropReportReply]()
	d.device.WhenComplete(func(dev *Device, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			return
		}
		var baseMsg = msg.NewBaseMsg().WithMethod("device.prop.report")
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
				var v = PropReportReply{*m}
				ret.Complete(&v)
			}
		})

	})
	return ret
}
