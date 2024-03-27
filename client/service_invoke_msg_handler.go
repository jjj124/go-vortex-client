package client

import (
	"github.com/jjj124/go-vortex-client/msg"
)

func NewServiceInvokeMsgHandler() ReceivedMsgHandler {
	return func(msg *msg.ReceivedMsg, client AdapterClient) {
		defer func() {

		}()
		if msg.Method() == "device.service.invoke" {
			handleServiceInvokeMsg(msg, client)
		}
	}
}

func handleServiceInvokeMsg(msg *msg.ReceivedMsg, client AdapterClient) {
	var headers = msg.Header()
	var pid = headers["adapter_pid"].(string)
	var did = headers["adapter_did"].(string)
	var hardwareId = headers["adapter_hardware_id"].(string)
	var dev Device
	if headers["adapter_gw_id"] != nil {
		var gwid = headers["adapter_gw_id"].(string)
		dev = NewSubDevice(pid, did, hardwareId, gwid)
	} else {
		dev = NewDevice(pid, did, hardwareId)
	}
	var payload = msg.Payload()
	var replyPayload = make(map[string]any, 0)
	for k, v := range payload {
		var fc = client.Components().ServiceHandlers()[k]
		if fc != nil {
			func() {
				defer func() {
					if recover() != nil {
						var tmp = make(map[string]any, 1)
						tmp["success"] = false
						tmp["error"] = recover()
						replyPayload[k] = tmp
					}
				}()
				var params = v.(map[string]any)
				var suc, r, err = fc(dev, params, client)
				var tmp = make(map[string]any, 1)
				if suc {
					tmp["success"] = true
					tmp["result"] = r
				} else {
					tmp["success"] = false
					tmp["error"] = err
				}
				replyPayload[k] = tmp
			}()

		}
	}
	var deliveryMsg = msg.NewReply(replyPayload)
	client.Delivery(deliveryMsg)
}
