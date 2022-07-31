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

	CommandRecentContest     = "recent_contest"
	CommandUserContestRecord = "user_contest_record"

	BlankSpace = "--------------"
)

var (
	PlaformList      = []string{PlatformCf, PlatformAt, PlatformNc, PlatformLg, PlatformVj, PlatformLc, PlaformCc}
	MsgHandlerMap    map[string]DefaultHandler
	TimerTaskManager *TimerTask
)
