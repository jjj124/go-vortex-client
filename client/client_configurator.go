package client

type Configurator interface {
	HandleServiceWith(serviceName string, fc DeviceServiceInvokeHandler) Configurator
	AnswerDeviceConnectStateWith(fc func(device Device) bool) Configurator
}

type configurator struct {
	component DefaultAdapterComponent
}

func (c *configurator) HandleServiceWith(serviceName string, fc DeviceServiceInvokeHandler) Configurator {
	c.component.serviceHandlers[serviceName] = fc
	return c
}

func (c *configurator) AnswerDeviceConnectStateWith(fc func(device Device) bool) Configurator {
	c.component.deviceConnectStateAsserter = fc
	return c
}
