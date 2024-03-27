package msg

import mqtt "github.com/eclipse/paho.mqtt.golang"

type DeliveryMsg struct {
	BaseMsg
	mqttClient mqtt.Client
}

func (d *DeliveryMsg) MqttClient() mqtt.Client {
	return d.mqttClient
}

func NewDeliveryMsg(baseMsg *BaseMsg) *DeliveryMsg {
	return &DeliveryMsg{
		BaseMsg:    *baseMsg,
		mqttClient: nil,
	}
}
