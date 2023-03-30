package model

const (
	UrlRobotBase      = "http://localhost:5700"
	UrlSendGroupMsg   = "send_msg"
	UrlSendPrivateMsg = "send_private_msg"

	Content
	SelfId = "2061972845"

	PlatformCf = "codeforces"
	PlatformAt = "atcoder"
	PlatformNc = "nowcoder"
	PlatformLg = "luogu"
	PlatformVj = "vjudge"
	PlatformLc = "leetcode"
	PlaformCc  = "codechef"

	BlankSpace = "--------------"
)

var (
	PlaformList      = []string{PlatformCf, PlatformAt, PlatformNc, PlatformLg, PlatformVj, PlatformLc, PlaformCc}
	MsgHandlerMap    = make(map[string]*DefaultHandler)
	TimerTaskManager *TimerTask

	CmdGetRecentContest              *DefaultHandler
	CmdGetUserContestRecord          *DefaultHandler
	CmdGetTimerContest               *DefaultHandler
	CmdFetchTimerContest             *DefaultHandler
	CmdAddTimerGroupId               *DefaultHandler
	CmdHelp                          *DefaultHandler
	CmdBindUser                      *DefaultHandler
	CmdQueryUser                     *DefaultHandler
	CmdFetchTimerPlatformUserRecord  *DefaultHandler
	CmdGetBindUserRecord             *DefaultHandler
	CmdFetchTimerBindUserDailyRecord *DefaultHandler
)
