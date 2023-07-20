package notify

import (
	"context"
)

type INotify interface {
	NotifyText(ctx context.Context, content string, atUserIds []string) error
	NotifyMarkdown(ctx context.Context, content string, atUserIds []string) error
	SetWebHook(ctx context.Context, webhook string)
	Name() string
	Webhook() string
}

func NewNotify(webhook string) INotify {
	n := notifyXs{webhook: webhook}
	return &n
}
