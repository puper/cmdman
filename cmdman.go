package cmdman

func NewCMDMan(cmdRunnerCreater func(*CMDMan, string) CMDRunner) *CMDMan {
	cmdMan := &CMDMan{
		cmdList:          NewCMDList(),
		cmdRunners:       map[string]CMDRunner{},
		cmdRunnerCreater: cmdRunnerCreater,
	}
	go cmdMan.mainLoop()
	return cmdMan
}

type CMDMan struct {
	cmdList          *CMDList
	cmdRunners       map[string]CMDRunner
	cmdRunnerCreater func(*CMDMan, string) CMDRunner
}

func (this *CMDMan) SendCMD(cmd *CMD) {
	this.cmdList.Put(cmd)
}

func (this *CMDMan) mainLoop() {
	for {
		select {
		case <-this.cmdList.CMDNotify():
			for {
				cmd := this.cmdList.Get()
				if cmd == nil {
					break
				}
				if cmd.Type == CMDTypeClose {
					if cmdRunner, ok := this.cmdRunners[cmd.Owner]; ok {
						if cmdRunner.Len() == 0 {
							cmd.SetResult(true)
							delete(this.cmdRunners, cmd.Owner)
						} else {
							cmd.SetResult(false)
						}
					} else {
						cmd.SetResult(true)
						delete(this.cmdRunners, cmd.Owner)
					}
				} else {
					cmdRunner, ok := this.cmdRunners[cmd.Owner]
					if !ok {
						cmdRunner = this.cmdRunnerCreater(this, cmd.Owner)
						this.cmdRunners[cmd.Owner] = cmdRunner
					}
					cmdRunner.SendCMD(cmd)
				}
			}
		}
	}
}
