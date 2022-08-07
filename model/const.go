package model

const (
	UrlRobotBase    = "http://localhost:5700"
	UrlSendGroupMsg = "send_msg"

	MethodSendGroupMsg = "POST"

	ContentTypeSendGroupMsg = "application/json"
	SelfId                  = "2061972845"

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

	CmdGetRecentContest     = DefaultHandler{Name: "get_recent_contest"}
	CmdGetUserContestRecord = DefaultHandler{Name: "get_user_contest_record"}
	CmdGetTimerContest      = DefaultHandler{Name: "get_timer_contest"}
	CmdFetchTimerContest    = DefaultHandler{Name: "fetch_timer_contest"}
	CmdAddTimerGroupId      = DefaultHandler{Name: "add_timer_group_id"}
	CmdHelp                 = DefaultHandler{Name: "help"}
)
