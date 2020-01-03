package cmdman

import (
	"encoding/json"
	"log"
)

type CMDRunner interface {
	SendCMD(*CMD)
	Len() int
}

type CMDRunnerBase struct {
	cmdMan    *CMDMan
	Owner     string
	cmdList   *CMDList
	handleCMD func(*CMD)
}

func NewCMDRunnerBase(cmdMan *CMDMan, owner string, handleCMD func(*CMD)) *CMDRunnerBase {
	cmdRunner := &CMDRunnerBase{
		cmdMan:    cmdMan,
		Owner:     owner,
		cmdList:   NewCMDList(),
		handleCMD: handleCMD,
	}
	if cmdRunner.handleCMD == nil {
		cmdRunner.handleCMD = func(cmd *CMD) {
			b, _ := json.Marshal(cmd.Params)
			log.Println("handleCMD: ", string(b))
		}
	}
	go cmdRunner.mainLoop()
	return cmdRunner
}

func (this *CMDRunnerBase) SendCMD(cmd *CMD) {
	this.cmdList.Put(cmd)
}

func (this *CMDRunnerBase) Len() int {
	return this.cmdList.list.Len()
}

func (this *CMDRunnerBase) mainLoop() {
	for {
		select {
		case <-this.cmdList.CMDNotify():
			for {
				cmd := this.cmdList.Get()
				if cmd == nil {
					break
				}
				this.handleCMD(cmd)
			}
		default:
			cmd := NewCMD(this.Owner, CMDTypeClose, nil)
			this.cmdMan.SendCMD(cmd)
			result := cmd.GetResult()
			if result.(bool) {
				return
			}
		}
	}
}
