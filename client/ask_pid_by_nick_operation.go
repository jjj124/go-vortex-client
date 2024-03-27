package client

import (
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	futures "github.com/jjj124/go-future"
)

type AskPidByNickReply struct {
	*msg.ReceivedMsg
}

func (d *AskPidByNickReply) GetPid() string {
	var x = d.Payload()["pid"]
	return x.(string)
}

type AskPidByNickOperation interface {
	WithNick(nick string) AskPidByNickOperation
	Execute() futures.Future[AskPidByNickReply]
}

type askPidByNickOperation struct {
	adapter AdapterClient
	nick    string
}

func (a *askPidByNickOperation) WithNick(nick string) AskPidByNickOperation {
	a.nick = nick
	return a
}

func (a *askPidByNickOperation) Execute() futures.Future[AskPidByNickReply] {
	var baseMsg = msg.NewBaseMsg()
	baseMsg.WithMethod("vortex.adapter.ask-pid-by-nick")
	var payload = baseMsg.Payload()
	payload["nick"] = a.nick
	var deliveryMsg = msg.NewDeliveryMsg(baseMsg)
	var f = a.adapter.Delivery(deliveryMsg)
	return futures.Then(f, func(m *msg.ReceivedMsg) *AskPidByNickReply {
		var v = AskPidByNickReply{m}
		a.adapter.Caches().CachePidNick(a.nick, v.GetPid())
		return &v
	})
}
func NewAskPidByNickOperation(adapter AdapterClient) AskPidByNickOperation {
	return &askPidByNickOperation{adapter: adapter}
}
