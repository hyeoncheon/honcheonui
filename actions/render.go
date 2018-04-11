package actions

import (
	"html/template"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/markbates/inflect"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
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
					return template.HTML(`<i class="fa fa-battery-full mixin-green"></i>`)
				case "status-true-false":
					return template.HTML(`<i class="fa fa-battery-full"></i>`)
				case "status-false-true":
					return template.HTML(`<i class="fa fa-battery-empty mixin-green"></i>`)
				case "status-false-false":
					return template.HTML(`<i class="fa fa-battery-empty"></i>`)
				default:
					return template.HTML(`<i class="fa fa-` + s + `"></i>`)
				}
			},
		},
	})
}
