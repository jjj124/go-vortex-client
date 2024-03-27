package client

import (
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	futures "github.com/jjj124/go-future"
)

type CreateDeviceReply struct {
	msg.ReceivedMsg
}

func (d *CreateDeviceReply) GetDid() string {
	var x = d.Payload()["did"]
	return x.(string)
}
func (d *CreateDeviceReply) GetPid() string {
	var x = d.Payload()["pid"]
	return x.(string)
}
func (d *CreateDeviceReply) GetHardwareId() string {
	var x = d.Payload()["hardware_id"]
	return x.(string)
}
func (d *CreateDeviceReply) GetGwId() string {
	var x = d.Payload()["gw_id"]
	if x == nil {
		return ""
	}
	return x.(string)
}
func (d *CreateDeviceReply) Wrap() Device {
	var devImpl Device
	if d.GetGwId() == "" {
		devImpl = NewDevice(d.GetPid(), d.GetDid(), d.GetHardwareId())
	} else {
		devImpl = NewSubDevice(d.GetPid(), d.GetDid(), d.GetHardwareId(), d.GetGwId())
	}
	return devImpl
}

type CreateDeviceOperation interface {
	WithPid(pid string) CreateDeviceOperation
	WithHardwareId(hardwareId string) CreateDeviceOperation
	WithGwId(gwId string) CreateDeviceOperation
	WithCreateIfNotExist(b bool) CreateDeviceOperation
	WithName(name string) CreateDeviceOperation
	Execute() futures.Future[CreateDeviceReply]
}

type defaultCreateDeviceOperation struct {
	adapter          AdapterClient
	pid              string
	hardwareId       string
	gwId             string
	name             string
	createIfNotExist bool
}

func (c *defaultCreateDeviceOperation) WithPid(pid string) CreateDeviceOperation {
	c.pid = pid
	return c
}

func (c *defaultCreateDeviceOperation) WithHardwareId(hardwareId string) CreateDeviceOperation {
	c.hardwareId = hardwareId
	return c
}

func (c *defaultCreateDeviceOperation) WithGwId(gwId string) CreateDeviceOperation {
	c.gwId = gwId
	return c
}

func (c *defaultCreateDeviceOperation) WithCreateIfNotExist(b bool) CreateDeviceOperation {
	c.createIfNotExist = b
	return c
}

func (c *defaultCreateDeviceOperation) WithName(name string) CreateDeviceOperation {
	c.name = name
	return c
}

func (c *defaultCreateDeviceOperation) Execute() futures.Future[CreateDeviceReply] {
	var baseMsg = msg.NewBaseMsg()
	baseMsg.WithMethod("vortex.adapter.device.create")
	var payload = baseMsg.Payload()
	payload["pid"] = c.pid
	if c.gwId != "" {
		payload["gw_id"] = c.gwId
	}
	if c.name != "" {
		payload["name"] = c.name
	}
	payload["create_if_not_exist"] = c.createIfNotExist
	payload["hardware_id"] = c.hardwareId
	var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
	var f = c.adapter.Delivery(deliveryMsg)
	return futures.Then(f, func(t *msg.ReceivedMsg) *CreateDeviceReply {
		var v = CreateDeviceReply{*t}
		var dev = v.Wrap()
		c.adapter.Caches().GetDeviceCacheByPid(dev.Pid()).PutDevice(dev)
		return &v
	})
}

func NewCreateDeviceOperation(adapter AdapterClient) CreateDeviceOperation {
	return &defaultCreateDeviceOperation{
		adapter:          adapter,
		createIfNotExist: true,
	}
}
