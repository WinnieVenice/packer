package model

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Task struct {
	Time string
	Task func()
}
type TimerTask struct {
	Cron     *cron.Cron
	TaskList []Task
}

func NewTimerTask() *TimerTask {
	tt := TimerTask{}
	return &tt
}

func GetTimerTask() *TimerTask {
	if TimerTaskManager == nil {
		TimerTaskManager = NewTimerTask()
	}
	return TimerTaskManager
}

func (tt *TimerTask) GetCron() *cron.Cron {
	if tt.Cron == nil {
		shc, _ := time.LoadLocation("Asia/Shanghai")
		tt.Cron = cron.New(cron.WithSeconds(), cron.WithLocation(shc))
	}
	return tt.Cron
}

func (tt *TimerTask) Add(time string, task func()) error {
	_, err := tt.GetCron().AddFunc(time, task)
	if err != nil {
		return err
	}
	return nil
}
func (tt *TimerTask) Run() {
	c := tt.GetCron()
	for _, task := range tt.TaskList {
		err := tt.Add(task.Time, task.Task)
		if err != nil {
			panic(err)
		}
	}
	// c.Start()
	c.Run()
}
