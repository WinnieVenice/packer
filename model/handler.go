package model

import (
	"context"
)

type MsgHandler func(ctx context.Context, commCtx map[string]string, param []string) error
type DefaultHandler func(commCtx map[string]string, param []string) error
