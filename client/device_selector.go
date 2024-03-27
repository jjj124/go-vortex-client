package client

import (
	futures "github.com/jjj124/go-future"
)

type DeviceSelector interface {
	WithNameIfCreate(name string) DeviceSelector
	WithCreateIfNotExist(b bool) DeviceSelector
	Get() futures.Future[Device]
	Op() DeviceOperations
}
type GatewayDeviceSelector interface {
	WithNameIfCreate(name string) GatewayDeviceSelector
	WithCreateIfNotExist(b bool) GatewayDeviceSelector
	Get() futures.Future[Device]
	Op() DeviceOperations
	SubDeviceOf(pidNick string, hardwareId string) DeviceSelector
}

type deviceSelector struct {
	adapter          AdapterClient
	pid              futures.Future[string]
	hardwareId       string
	name             string
	gw               futures.Future[Device]
	createIfNotExist bool
}

func (d *deviceSelector) Get() futures.Future[Device] {
	var devFuture = futures.NewFuture[Device]()
	if d.gw != nil {
		var unionF = futures.And(d.pid, d.gw)
		unionF.WhenComplete(func(s *futures.Tuple2[string, Device], err error) {
			if err != nil {
				devFuture.CompleteExceptionally(err)
				return
			}
			var pid, _ = d.pid.BlockingGet()
			var gw, _ = d.gw.BlockingGet()
			var hardwareId = (*gw).HardwareId() + "@@@" + d.hardwareId
			var dev = d.adapter.Caches().GetDeviceCacheByPid(*pid).GetDeviceByHardwareId(hardwareId)
			if dev != nil {
				devFuture.Complete(&dev)
				return
			}
			var createDevFuture = d.adapter.AdapterOperation().CreateDeviceOperation().WithPid(*pid).WithHardwareId(d.hardwareId).WithCreateIfNotExist(d.createIfNotExist).WithName(d.name).WithGwId((*gw).Did()).Execute()
			createDevFuture.WhenComplete(func(c *CreateDeviceReply, err error) {
				if err != nil {
					devFuture.CompleteExceptionally(err)
					return
				}
				var devImpl = c.Wrap()
				devFuture.Complete(&devImpl)
			})
		})
	} else {
		d.pid.WhenComplete(func(s *string, err error) {
			if err != nil {
				devFuture.CompleteExceptionally(err)
				return
			}
			var pid = *s
			var dev = d.adapter.Caches().GetDeviceCacheByPid(pid).GetDeviceByHardwareId(d.hardwareId)
			if dev != nil {
				devFuture.Complete(&dev)
				return
			}
			var createDevFuture = d.adapter.AdapterOperation().CreateDeviceOperation().WithPid(pid).WithHardwareId(d.hardwareId).WithCreateIfNotExist(d.createIfNotExist).WithName(d.name).Execute()
			createDevFuture.WhenComplete(func(c *CreateDeviceReply, err error) {
				if err != nil {
					devFuture.CompleteExceptionally(err)
					return
				}
				var dev = NewDevice(c.GetPid(), c.GetDid(), c.GetHardwareId())
				devFuture.Complete(&dev)
			})

		})
	}
	return devFuture

}

func (d *deviceSelector) WithCreateIfNotExist(b bool) DeviceSelector {
	d.createIfNotExist = b
	return d
}

func (d *deviceSelector) WithNameIfCreate(name string) DeviceSelector {
	d.name = name
	return d
}

func (d *deviceSelector) Op() DeviceOperations {
	return NewDeviceOperations(d.adapter, d.Get())
}

func MakeDeviceOf(adapter AdapterClient, pid string, hardwareId string) DeviceSelector {
	var ret = &deviceSelector{
		adapter:          adapter,
		pid:              futures.Just[string](&pid),
		hardwareId:       hardwareId,
		name:             "",
		gw:               nil,
		createIfNotExist: true,
	}
	return ret
}

type gatewayDeviceSelector struct {
	deviceSelector
}

func (d *gatewayDeviceSelector) WithNameIfCreate(name string) GatewayDeviceSelector {
	d.name = name
	return d
}

func (d *gatewayDeviceSelector) WithCreateIfNotExist(b bool) GatewayDeviceSelector {
	d.createIfNotExist = b
	return d
}

func (d *gatewayDeviceSelector) SubDeviceOf(pidNick string, hardwareId string) DeviceSelector {

	var pid = d.adapter.Caches().GetPidByNick(pidNick)
	var pidFuture futures.Future[string]
	if pid != "" {
		pidFuture = futures.Just(&pid)
	} else {
		pidFuture = futures.Then(d.adapter.AdapterOperation().AskPidByNickOperation().WithNick(pidNick).Execute(), func(t *AskPidByNickReply) *string {
			var pid = t.GetPid()
			return &pid
		})
	}
	return &deviceSelector{
		adapter:          d.adapter,
		pid:              pidFuture,
		hardwareId:       hardwareId,
		name:             "",
		gw:               d.Get(),
		createIfNotExist: true,
	}
}

func NewGatewayOf(adapter AdapterClient, pid string, hardwareId string) GatewayDeviceSelector {
	return &gatewayDeviceSelector{
		deviceSelector: deviceSelector{
			adapter:          adapter,
			pid:              futures.Just(&pid),
			hardwareId:       hardwareId,
			createIfNotExist: true,
		},
	}
}
