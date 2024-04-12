package client

import "github.com/jjj124/go-vortex-client/msg"

func NewAskRecentMsgHandler() ReceivedMsgHandler {
	return func(msg *msg.ReceivedMsg, client AdapterClient) {
		defer func() {

		}()
		if msg.Method() == "vortex.adapter.ask-recent-msg" {
			handleAskRecentMsg(msg, client)
		}
	}
}

func handleAskRecentMsg(msg *msg.ReceivedMsg, client AdapterClient) {
	var payload = msg.Payload()
	var recentMsg RecentMsg
	if payload["type"] == "send" {
		recentMsg = client.Components().RecentSendMsg()
	} else {
		recentMsg = client.Components().RecentReceivedMsg()
	}
	var snapshot = recentMsg.Snapshot()
	var list = make([]any, len(snapshot))
	for index, item := range snapshot {
		list[index] = item.ToMap()
	}
	var replyPayload = make(map[string]any)
	replyPayload["list"] = list
	client.Delivery(msg.NewReply(replyPayload))
}
