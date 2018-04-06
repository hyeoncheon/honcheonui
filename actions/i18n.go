package actions

import (
	"github.com/gobuffalo/buffalo"
)

func t(c buffalo.Context, str string, args ...interface{}) string {
	s := T.Translate(c, str, args...)
	if s == str {
		c.Logger().WithField("category", "i18n").Warnf("UNTRANSLATED: %v", str)
	}
	return s
}
