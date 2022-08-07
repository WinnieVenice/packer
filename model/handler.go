package model

import (
	"context"
	"fmt"
)

type MsgHandler func(ctx context.Context, commCtx map[string]string, param []string) error

// type DefaultHandler func(commCtx map[string]string, param []string) error

type HttpHandler func(commCtx map[string]string, param []string) error

type DefaultHandler struct {
	Name    string
	Content string
	Handler HttpHandler

	CmdMapp []string
}

func (dh *DefaultHandler) String() string {
	return fmt.Sprintf("%s\n命令: %+v", dh.Content, dh.CmdMapp)
}
