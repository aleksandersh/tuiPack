package application

import (
	"github.com/rivo/tview"
)

type UiController struct {
	app *tview.Application
}

func NewController(app *tview.Application) *UiController {
	return &UiController{app: app}
}

func (controller UiController) Close() {
	if controller.app == nil {
		return
	}
	controller.app.Stop()
}
