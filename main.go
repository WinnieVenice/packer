package main

import (
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/service/cq"
)

func main() {

	go cq.Run()
	go model.GetTimerTask().Run()

	select {}
}
