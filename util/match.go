package util

import (
	"fmt"
	"packer/model"
)

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
	matchCommandMap = make(map[string]string)
)

func AddMatchCommand(cmd string, handler *model.DefaultHandler) bool {
	if len(cmd) <= 0 || handler == nil {
		fmt.Printf("AddMatchCommand failed, param invaild, cmd = (%+v), handlers = (%+v)\n", cmd, handler)
		return false
	}
	if _, ok := matchCommandMap[cmd]; ok {
		fmt.Printf("AddMatchCommand failed, cmd existed, cmd = (%+v)\n", cmd)
		return false
	}

	matchCommandMap[cmd] = handler.Name
	handler.CmdMapp = append(handler.CmdMapp, cmd)

	fmt.Printf("AddMatchCommand succ, handler = (%+v), cmdMapp = (%+v)\n", handler.Name, handler.CmdMapp)
	return true
}

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
