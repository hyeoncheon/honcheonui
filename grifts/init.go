package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hyeoncheon/honcheonui/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
