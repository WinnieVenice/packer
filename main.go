package main

import (
	"packer/handlers"
	"packer/http"
	"packer/model"
)

func init() {
	model.MsgHandlerMap = map[string]model.DefaultHandler{
		"recent_contest":             handlers.HandlerWrapper("GetRecentContest", handlers.GetRecentContest),
		"user_contest_record":        handlers.HandlerWrapper("GetUserContestRecord", handlers.GetUserContestRecord),
		"get_timer_contest":          handlers.HandlerWrapper("GetTimerContest", handlers.GetTimerAllContest),
		"fetch_timer_recent_contest": handlers.HandlerWrapper("FetchTimerRecentContest", handlers.FetchTimerRecentContest),
		"add_timer_group_id":         handlers.HandlerWrapper("AddTimerGroupId", handlers.AddTimerGroupId),
		"help":                       handlers.HandlerWrapper("Help", handlers.GetAllCommand),
	}
	model.GetTimerTask().TaskList = []model.Task{
		{
			Time: "0 0 4 * * *", // 对齐原神每日更新时间
			Task: handlers.TimerRecentContest,
		},
		{
			Time: "0 0 10 * * *",
			Task: handlers.TimerDailyQuestion,
		},
	}
}

func main() {

	go http.Run()
	go model.GetTimerTask().Run()

	select {}
}
