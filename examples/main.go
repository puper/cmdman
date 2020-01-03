package main

import (
	"log"
	"time"

	"github.com/puper/cmdman"
)

var (
	a = []interface{}{}
)

func main() {
	cmdMan := cmdman.NewCMDMan(func(cmdMan *cmdman.CMDMan, owner string) cmdman.CMDRunner {
		return cmdman.NewCMDRunnerBase(cmdMan, owner, func(cmd *cmdman.CMD) {
			log.Println("start: ", cmd.Params)
			a = append(a, cmd.Params)
			time.Sleep(time.Millisecond * 25)
			log.Println("end: ", cmd.Params)
		})
	})
	Run(cmdMan, "test", 0, 100)
	Run(cmdMan, "test", 100, 200)
	Run(cmdMan, "test", 300, 400)
	Run(cmdMan, "test", 1000, 1100)
	Run(cmdMan, "test", 800, 900)
	Run(cmdMan, "test", 600, 700)
	time.Sleep(5 * time.Second)
}

func Run(cmdMan *cmdman.CMDMan, owner string, from, to int) {
	go func() {
		for i := from; i < to; i++ {
			cmdMan.SendCMD(cmdman.NewCMD(owner, "", i))
			time.Sleep(time.Millisecond * 50)
		}
	}()
}
