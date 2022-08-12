package qqbot

import (
	"fmt"

	"github.com/WinnieVenice/packer/core/timer"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/util"
)

func init() {
	model.MsgHandlerMap = make(map[string]*model.DefaultHandler)
	model.CmdGetRecentContest = HandlerWrapper(map[string]string{
		"name":    "get_recent_contest",
		"content": "功能: 拉取某平台最近比赛, 用法: 命令 平台",
	}, SendRecentContest)
	model.MsgHandlerMap[model.CmdGetRecentContest.Name] = &model.CmdGetRecentContest

	model.CmdGetUserContestRecord = HandlerWrapper(map[string]string{
		"name":    "get_user_contest_record",
		"content": "功能: 拉取用户记录, 用法: 命令 平台 用户名",
	}, SendUserContestRecord)
	model.MsgHandlerMap[model.CmdGetUserContestRecord.Name] = &model.CmdGetUserContestRecord

	model.CmdGetTimerContest = HandlerWrapper(map[string]string{
		"name":    "get_timer_contest",
		"content": "功能: 拉取定时推送比赛列表, 用法: 命令",
	}, GetTimerAllContest)
	model.MsgHandlerMap[model.CmdGetTimerContest.Name] = &model.CmdGetTimerContest

	model.CmdFetchTimerContest = HandlerWrapper(map[string]string{
		"name":    "fetch_timer_contest",
		"content": "功能: 更新定时推送比赛列表, 用法: 命令",
	}, FetchTimerRecentContest)
	model.MsgHandlerMap[model.CmdFetchTimerContest.Name] = &model.CmdFetchTimerContest

	model.CmdAddTimerGroupId = HandlerWrapper(map[string]string{
		"name":    "add_timer_group_id",
		"content": "功能: 将本群加入定时推送比赛列表, 用法: 命令",
	}, AddTimerGroupId)
	model.MsgHandlerMap[model.CmdAddTimerGroupId.Name] = &model.CmdAddTimerGroupId

	model.CmdHelp = HandlerWrapper(map[string]string{
		"name":    "help",
		"content": "功能: 帮助, 用法: 命令",
	}, SendHelp)
	model.MsgHandlerMap[model.CmdHelp.Name] = &model.CmdHelp

	fmt.Printf("MsgHandlerMap = (%+v)\n", model.MsgHandlerMap)

	util.AddMatchCommand("比赛", &model.CmdGetRecentContest)
	util.AddMatchCommand("康康", &model.CmdGetUserContestRecord)
	util.AddMatchCommand("推送列表", &model.CmdGetTimerContest)
	util.AddMatchCommand("更新推送", &model.CmdFetchTimerContest)
	util.AddMatchCommand("加入推送", &model.CmdAddTimerGroupId)
	util.AddMatchCommand("帮助", &model.CmdHelp)

	model.GetTimerTask().TaskList = []model.Task{
		{
			Time: "0 0 4 * * *", // 对齐原神每日更新时间
			Task: timer.RecentContest,
		},
		{
			Time: "0 0 10 * * *",
			Task: timer.DailyQuestion,
		},
	}
}
