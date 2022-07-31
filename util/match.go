package util

import "packer/model"

var (
	matchPlatformMap = map[string]string{
		"cf": model.PlatformCf,
		"at": model.PlatformAt,
		"nc": model.PlatformNc,
		"lg": model.PlatformLg,
		"vj": model.PlatformVj,
		"lc": model.PlatformLc,
		"cc": model.PlaformCc,
	}
	matchCommandMap = map[string]string{
		"比赛":   model.CommandRecentContest,
		"最近比赛": model.CommandRecentContest,
		"康康":   model.CommandUserContestRecord,
		"让我康康": model.CommandUserContestRecord,
		"检查身体": model.CommandUserContestRecord,
	}
)

func MatchPlatform(platform string) string {
	if _, ok := matchPlatformMap[platform]; !ok {
		return platform
	}
	return matchPlatformMap[platform]
}

func MatchCommand(cmd string) string {
	if _, ok := matchCommandMap[cmd]; !ok {
		return cmd
	}
	return matchCommandMap[cmd]
}
