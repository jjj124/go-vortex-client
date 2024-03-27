package client

import (
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-vortex-client/msg"
	"github.com/spf13/viper"
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
