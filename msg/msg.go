package msg

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jjj124/go-vortex-client/utils"
	"time"
)

type Msg interface {
	Method() string
	Ts() int64
	MsgId() string
	Payload() map[string]any
	Header() map[string]any
	Error() map[string]any
	Marshal() ([]byte, error)
	ToString() string
}

type MqttClientBound interface {
	MqttClient() mqtt.Client
}

type BaseMsg struct {
	ts      int64
	method  string
	msgId   string
	payload map[string]any
	header  map[string]any
	error   map[string]any
}

func (b *BaseMsg) Error() map[string]any {
	return b.error
}

func (b *BaseMsg) ToString() string {
	var bytes, err = json.Marshal(b.toMap())
	if err != nil {
		panic(err)
	}
	var str = string(bytes)
	return str
}
func (b *BaseMsg) toMap() map[string]any {
	var j = make(map[string]any, 5)
	j["ts"] = b.ts
	j["method"] = b.method
	j["payload"] = b.payload
	j["msg_id"] = b.msgId
	j["header"] = b.header
	j["error"] = b.error
	return j
}

func (b *BaseMsg) Marshal() ([]byte, error) {
	var j = b.toMap()
	return json.Marshal(j)
}

func (b *BaseMsg) Method() string {
	return b.method
}

func (b *BaseMsg) Ts() int64 {
	return b.ts
}

func (b *BaseMsg) MsgId() string {
	return b.msgId
}

func (b *BaseMsg) Payload() map[string]any {
	return b.payload
}

func (b *BaseMsg) Header() map[string]any {
	return b.header
}

func (b *BaseMsg) WithMethod(method string) *BaseMsg {
	b.method = method
	return b
}

func NewBaseMsg() *BaseMsg {
	var ts = time.Now().UnixMilli()
	var msgId = utils.RandString(10)
	var payload = make(map[string]any)
	var header = make(map[string]any)
	return &BaseMsg{
		ts,
		"",
		msgId,
		payload,
		header,
		nil,
	}
}
