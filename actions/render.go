package actions

import (
	"html/template"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"github.com/gofrs/uuid"
	"github.com/markbates/inflect"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public")

func init() {
	r = render.New(render.Options{
		HTMLLayout:   "application.html",
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		Helpers: render.Helpers{
			"titleize": inflect.Titleize,
			"iconize": func(s string) template.HTML {
				switch s {
				case "vm":
					return template.HTML(`<i class="fa fa-paper-plane"></i>`)
				case "bm":
					return template.HTML(`<i class="fa fa-plane"></i>`)
				case "app":
					return template.HTML(`<i class="fa fa-cog"></i>`)
				case "database":
					return template.HTML(`<i class="fa fa-database"></i>`)
				case "status-true-true":
					return template.HTML(`<i class="fa fa-play-circle"></i>`)
				case "status-true-false":
					return template.HTML(`<i class="fa fa-power-off mixin-orange"></i>`)
				case "status-false-true":
					return template.HTML(`<i class="fa fa-signal"></i>`)
				case "status-false-false":
					return template.HTML(`<i class="fa fa-power-off mixin-red"></i>`)
				default:
					return template.HTML(`<i class="fa fa-` + s + `"></i>`)
				}
			},
			"has": func(a []string, v string) bool {
				for _, e := range a {
					if e == v {
						return true
					}
				}
				return false
			},
			"uuidTrucate": func(u uuid.UUID, l ...int) string {
				if len(l) > 0 {
					return u.String()[0:l[0]]
				}
				return u.String()
			},
		},
	})
}
