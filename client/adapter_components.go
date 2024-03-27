package client

import "github.com/jjj124/go-metrics"

type AdapterComponents interface {
	RecentSendMsg() RecentMsg
	RecentReceivedMsg() RecentMsg
	DeviceConnectStateAsserter() func(device Device) bool
	ServiceHandlers() map[string]DeviceServiceInvokeHandler
	MetricsRegistry() metrics.Registry
}

type DefaultAdapterComponent struct {
	recentSendMsg              RecentMsg
	recentReceivedMsg          RecentMsg
	deviceConnectStateAsserter func(device Device) bool
	serviceHandlers            map[string]DeviceServiceInvokeHandler
	metricsRegistry            metrics.Registry
}

func (a *DefaultAdapterComponent) MetricsRegistry() metrics.Registry {
	return a.metricsRegistry
}

func (a *DefaultAdapterComponent) ServiceHandlers() map[string]DeviceServiceInvokeHandler {
	return a.serviceHandlers
}

func (a *DefaultAdapterComponent) RecentSendMsg() RecentMsg {
	return a.recentSendMsg
}

func (a *DefaultAdapterComponent) RecentReceivedMsg() RecentMsg {
	return a.recentReceivedMsg
}

func (a *DefaultAdapterComponent) DeviceConnectStateAsserter() func(device Device) bool {
	return a.deviceConnectStateAsserter
}

func NewAdapterComponent() DefaultAdapterComponent {
	return DefaultAdapterComponent{
		recentReceivedMsg: NewRecentMsg(),
		recentSendMsg:     NewRecentMsg(),
		deviceConnectStateAsserter: func(device Device) bool {
			return true
		},
		serviceHandlers: make(map[string]DeviceServiceInvokeHandler),
		metricsRegistry: metrics.NewRegistry(),
	}
}
