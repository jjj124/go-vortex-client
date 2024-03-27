package client

import (
	"container/list"
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	"sync"
)

type RecentMsg interface {
	Push(msg msg.Msg)
}

type recentMsg struct {
	lock  sync.Locker
	queue *list.List
	max   int
}

func (a *recentMsg) Push(msg msg.Msg) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.queue.PushBack(msg)
	if a.queue.Len() > a.max {
		var first = a.queue.Front()
		a.queue.Remove(first)
	}
}

func (a *recentMsg) Snapshot() []msg.Msg {
	a.lock.Lock()
	defer a.lock.Unlock()
	var ret = make([]msg.Msg, a.queue.Len())
	var i = 0
	for e := a.queue.Front(); e != nil; e = e.Next() {
		ret[i] = e.Value.(msg.Msg)
		i++
	}
	return ret
}

func NewRecentMsg() RecentMsg {
	return &recentMsg{
		lock:  &sync.Mutex{},
		queue: list.New(),
		max:   30,
	}
}
