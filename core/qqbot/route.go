package qqbot

import (
	"fmt"

	"github.com/WinnieVenice/packer/core/timer"
	"github.com/WinnieVenice/packer/model"
	"github.com/WinnieVenice/packer/util"
)

func init() {
	model.MsgHandlerMap = make(map[string]*model.DefaultHandler)

	model.CmdGetRecentContest = register(SendRecentContest, map[string]string{
		"name":    "get_recent_contest",
		"content": "功能: 拉取某平台最近比赛, 用法: 命令 平台",
	})

	model.CmdGetUserContestRecord = register(SendUserContestRecord, map[string]string{
		"name":    "get_user_contest_record",
		"content": "功能: 拉取用户记录, 用法: 命令 平台 用户名",
	})

	model.CmdGetTimerContest = register(GetTimerAllContest, map[string]string{
		"name":    "get_timer_contest",
		"content": "功能: 拉取定时推送比赛列表, 用法: 命令",
	})

	model.CmdFetchTimerContest = register(FetchTimerRecentContest, map[string]string{
		"name":    "fetch_timer_contest",
		"content": "功能: 更新定时推送比赛列表, 用法: 命令",
	})

	model.CmdAddTimerGroupId = register(AddTimerGroupId, map[string]string{
		"name":    "add_timer_group_id",
		"content": "功能: 将本群加入定时推送比赛列表, 用法: 命令",
	})

	model.CmdHelp = register(SendHelp, map[string]string{
		"name":    "help",
		"content": "功能: 帮助, 用法: 命令",
	})

	model.CmdBindUser = register(BindUser, map[string]string{
		"name":    "bind_user",
		"content": "功能: 绑定当前qq与oj平台id, 用法: 命令 平台1:id1 平台2:id2 ...",
	})

	model.CmdQueryUser = register(QueryUser, map[string]string{
		"name":    "query_user",
		"content": "功能: 查询当前qq的绑定信息, 用法: 命令",
	})

	model.CmdFetchTimerPlatformUserRecord = register(FetchTimerPlatformUserRecord, map[string]string{
		"name":    "fetch_timer_platform_user_record",
		"content": "功能: 更新定时用户记录, 用法: 命令",
	})

	model.CmdGetBindUserRecord = register(GetBindUserRecord, map[string]string{
		"name":    "get_bind_user_record",
		"content": "功能: 获取绑定oj的记录, 用法: 命令",
	})

	model.CmdFetchTimerBindUserDailyRecord = register(FetchTimerBindUserDailyRecord, map[string]string{
		"name":    "fetch_timer_bind_user_daily_record",
		"content": "功能: 更新每日训练记录, 用法: 命令",
	})

	fmt.Printf("MsgHandlerMap = (%+v)\n", model.MsgHandlerMap)

	util.AddMatchCommand(model.CmdGetRecentContest.Name, model.CmdGetRecentContest)
	util.AddMatchCommand("比赛", model.CmdGetRecentContest)

	util.AddMatchCommand(model.CmdGetUserContestRecord.Name, model.CmdGetUserContestRecord)
	util.AddMatchCommand("康康", model.CmdGetUserContestRecord)

	util.AddMatchCommand(model.CmdGetTimerContest.Name, model.CmdGetTimerContest)
	util.AddMatchCommand("推送列表", model.CmdGetTimerContest)

	util.AddMatchCommand(model.CmdFetchTimerContest.Name, model.CmdFetchTimerContest)
	util.AddMatchCommand("更新推送", model.CmdFetchTimerContest)

	util.AddMatchCommand(model.CmdAddTimerGroupId.Name, model.CmdAddTimerGroupId)
	util.AddMatchCommand("加入推送", model.CmdAddTimerGroupId)

	util.AddMatchCommand(model.CmdHelp.Name, model.CmdHelp)
	util.AddMatchCommand("帮助", model.CmdHelp)

	util.AddMatchCommand(model.CmdBindUser.Name, model.CmdBindUser)
	util.AddMatchCommand("绑定", model.CmdBindUser)

	util.AddMatchCommand(model.CmdQueryUser.Name, model.CmdQueryUser)
	util.AddMatchCommand("查询绑定", model.CmdQueryUser)

	util.AddMatchCommand(model.CmdFetchTimerPlatformUserRecord.Name, model.CmdFetchTimerPlatformUserRecord)
	util.AddMatchCommand("更新用户记录", model.CmdFetchTimerPlatformUserRecord)

	util.AddMatchCommand(model.CmdGetBindUserRecord.Name, model.CmdGetBindUserRecord)
	util.AddMatchCommand("查询绑定oj记录", model.CmdGetBindUserRecord)

	util.AddMatchCommand(model.CmdFetchTimerBindUserDailyRecord.Name, model.CmdFetchTimerBindUserDailyRecord)
	util.AddMatchCommand("更新训练记录", model.CmdFetchTimerBindUserDailyRecord)

	model.GetTimerTask().TaskList = []model.Task{
		{
			Time: "0 0/1 * * * *",
			Task: timer.RecentContest,
		},
		{
			Time: "0 0 10 * * *",
			Task: timer.DailyQuestion,
		},
		{
			Time: "0 0 10 * * *",
			Task: timer.TimerPlatformUserRecord,
		},
		{
			Time: "0 0 12 * * *",
			Task: timer.TimerBindUserDailyRecord,
		},
		{
			Time: "0 0 0/6 * * *",
			Task: timer.TimerJDSD,
		},
	}
}

func MustInit() {
}

func register(f model.MsgHandler, conf map[string]string) *model.DefaultHandler {
	handler := HandlerWrapper(conf, f)
	model.MsgHandlerMap[handler.Name] = &handler
	return &handler
}
