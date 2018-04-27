package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"

	"github.com/hyeoncheon/honcheonui/workers"
)

// AdminHandler render admin page.
func AdminHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("admin.html"))
}

// AdminSyncNotification queues notification sync job and redirect to admin root.
func AdminSyncNotification(c buffalo.Context) error {
	if err := workers.Run(workers.WorkerNotificationWatch, nil); err != nil {
		c.Flash().Add("danger", t(c, "Could.not.start.background.sync"))
	} else {
		c.Flash().Add("success", t(c, "Notifications.will.be.synced.in.background"))
	}

	return c.Redirect(http.StatusSeeOther, "/admin")
}
