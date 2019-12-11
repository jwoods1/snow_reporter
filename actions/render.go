package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"snow_watch/models"
	"time"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
			"timeFix":    timeFix,
			"lastReport": lastReport,
		},
	})
}
func timeFix(d time.Time) string {
	return d.Format("January 02, 2006")
}
func lastReport(reports []models.DailyReport) models.DailyReport {
	if len(reports) > 0 {
		return reports[len(reports)-1]

	}
	return models.DailyReport{}
}
