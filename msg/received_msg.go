package msg

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ReceivedMsg struct {
	BaseMsg
	mqttClient mqtt.Client
}

func (r *ReceivedMsg) NewReply(payload map[string]any) *DeliveryMsg {
	var base = NewBaseMsg()
	base.method = r.Method() + ".reply"
	base.payload = payload
	base.msgId = r.MsgId()
	if r.Header()["adapter_pid"] != nil {
		base.Header()["adapter_pid"] = r.Header()["adapter_pid"]
	}
	if r.Header()["adapter_did"] != nil {
		base.Header()["adapter_did"] = r.Header()["adapter_did"]
	}
	if r.Header()["adapter_hardware_id"] != nil {
		base.Header()["adapter_hardware_id"] = r.Header()["adapter_hardware_id"]
	}
	return NewDeliveryMsg(base)
}
func (r *ReceivedMsg) NewErrorReply(err error) {
	var base = NewBaseMsg()
	base.method = r.Method() + ".reply"
	base.error["msg"] = err
}

func (r *ReceivedMsg) MqttClient() mqtt.Client {
	return r.mqttClient
}

func NewReceivedMsg(bytes []byte, mqttClient mqtt.Client) (*ReceivedMsg, error) {
	var m = make(map[string]any, 0)

	var err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}
	var header map[string]any
	if m["header"] != nil {
		header = m["header"].(map[string]any)
	} else {
		header = make(map[string]any, 0)
	}
	var e map[string]any
	if m["error"] != nil {
		e = m["error"].(map[string]any)
	} else {
		e = nil
	}

	var baseMsg = BaseMsg{
		ts:      int64(m["ts"].(float64)),
		method:  m["method"].(string),
		msgId:   m["msg_id"].(string),
		payload: m["payload"].(map[string]any),
		header:  header,
		error:   e,
	}
	return &ReceivedMsg{
		BaseMsg:    baseMsg,
		mqttClient: mqttClient,
	}, nil
}
