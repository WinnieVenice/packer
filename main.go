package main

import (
	"github.com/WinnieVenice/packer/backend"
	"github.com/WinnieVenice/packer/core/qqbot"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/service/cq"
)

func main() {
	qqbot.MustInit()
	backend.MustInit()

	go cq.Run()
	go model.GetTimerTask().Run()

	select {}
}
