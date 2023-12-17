package application

type Application struct {
	Ui *UiController
}

func NewApplication(ui *UiController) *Application {
	return &Application{Ui: ui}
}
