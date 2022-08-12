package main

import (
	"github.com/WinnieVenice/packer/core/qqbot"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/service/cq"
)

func main() {
	qqbot.MustInit()

	go cq.Run()
	go model.GetTimerTask().Run()

	select {}
}
