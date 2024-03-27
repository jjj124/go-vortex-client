package client

import "github.com/Vortex-ECO/Vortex-SDk-GO/msg"

func NewAskDeviceConnectStateHandler() ReceivedMsgHandler {
	return func(msg *msg.ReceivedMsg, client AdapterClient) {
		defer func() {

		}()
		if msg.Method() == "vortex.adapter.ask-device-connect-state" {
			handleAskDeviceConnectHandler(msg, client)
		}
	}
}

func handleAskDeviceConnectHandler(msg *msg.ReceivedMsg, client AdapterClient) {
	var payload = msg.Payload()
	var dev Device = nil
	var pid = payload["pid"].(string)
	var did = payload["did"].(string)
	var hardwareId = payload["hardware_id"].(string)
	if payload["gw_id"] != nil {
		var gwId = payload["gw_id"].(string)
		dev = NewSubDevice(pid, did, hardwareId, gwId)
	} else {
		dev = NewDevice(pid, did, hardwareId)
	}
	var connected = client.Components().DeviceConnectStateAsserter()(dev)
	var reply = make(map[string]any, 1)
	reply["connected"] = connected
	client.Delivery(msg.NewReply(reply))
}
