package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
	"github.com/spf13/viper"
)

const Version = "0.0.1"

const (
	DebugLevelNone   = 0
	DebugLevelLow    = 1
	DebugLevelMedium = 2
	DebugLevelHigh   = 3
)
const (
	metricsAdapterMsgRecv      = "adapter.msg.recv"
	metricsAdapterMsgSend      = "adapter.msg.send"
	metricsAdapterModelMsgSend = "adapter.model.send"
)
const (
	DevicePropReport  = "device.prop.report"
	DeviceEventReport = "device.event.report"
)

type AdapterClient interface {
	Pid() string
	ClientId() string
	Start() futures.Future[string]
	Options() AdapterOptions
	Delivery(msg *msg.DeliveryMsg) futures.Future[msg.ReceivedMsg]
	AdapterOperation() AdapterOperations
	Caches() AdapterCaches
	DeviceOf(hardwareId string) DeviceSelector
	GwDeviceOf(hardwareId string) GatewayDeviceSelector
	Shutdown()
	Configurator() Configurator
	Components() AdapterComponents
	Viper() *viper.Viper
}
type ReceivedMsgHandler func(msg *msg.ReceivedMsg, client AdapterClient)

type DeviceServiceInvokeHandler func(device Device, params map[string]any, client AdapterClient) (bool, map[string]any, error)
