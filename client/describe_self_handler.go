package client

import (
	"github.com/jjj124/go-vortex-client/msg"
	"runtime"
)

func NewDescribeSelfHandler() ReceivedMsgHandler {
	return func(msg *msg.ReceivedMsg, client AdapterClient) {
		defer func() {

		}()
		if msg.Method() == "vortex.adapter.describe-self" {
			handleDescribeSelf(msg, client)
		}
	}
}

func handleDescribeSelf(msg *msg.ReceivedMsg, client AdapterClient) {
	var payload = make(map[string]any)
	var sdk = make(map[string]any)
	sdk["lang"] = "Golang"
	sdk["version"] = Version
	payload["sdk"] = sdk

	var metrics = client.Components().MetricsRegistry().GetAll()
	payload["metrics"] = metrics

	var options = make(map[string]any)
	options["pid"] = client.Options().Pid()
	options["connection_num"] = client.Options().ConnectNum()
	options["debug_level"] = client.Options().DebugLevel()
	options["client_id"] = client.Options().ClientId()
	payload["options"] = options

	var system = make(map[string]any)
	system["os"] = runtime.GOOS
	system["arch"] = runtime.GOARCH
	payload["system"] = system

	var reply = msg.NewReply(payload)
	client.Delivery(reply)
}
