package cmdman

import (
	"container/list"
	"sync"
)

func NewCMDList() *CMDList {
	return &CMDList{
		list:      list.New(),
		cmdNotify: make(chan struct{}, 1),
	}
}

type CMDList struct {
	mutex     sync.RWMutex
	list      *list.List
	cmdNotify chan struct{}
}

func (this *CMDList) Len() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.list.Len()
}

func (this *CMDList) Put(cmd *CMD) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if cmd != nil {
		this.list.PushBack(cmd)
		select {
		case this.cmdNotify <- struct{}{}:
		default:
		}
	}
}

func (this *CMDList) Get() *CMD {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	e := this.list.Front()
	if e != nil {
		this.list.Remove(e)
		return e.Value.(*CMD)
	}
	return nil
}

func (this *CMDList) CMDNotify() chan struct{} {
	return this.cmdNotify
}

func (this *CMDList) OnNewCMD(cb func(*CMD)) {
	for {
		select {
		case <-this.cmdNotify:
			for {
				cmd := this.Get()
				if cmd == nil {
					break
				}
				cb(cmd)
			}
		}
	}
}
