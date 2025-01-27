package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	app := fyneApp.New()
	mainWindow := setupMainWindow(app)

	setupSystrayMenu(app, mainWindow)

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	chromePath, foundChrome := launcher.LookPath()
	if foundChrome {
		fmt.Printf("chrome path = %q\n", chromePath)
	}

	mainWindow.Resize(fyne.NewSize(800, 600))

	app.Lifecycle().SetOnStarted(func() {
		if !foundChrome {
			noChrome := dialog.NewError(errors.New("Google Chrome (or equivalent) was not found on your computer.\n\nPlease install Google Chrome (or chromium)"), mainWindow)
			noChrome.Resize(fyne.NewSize(500, 200))
			noChrome.Show()
		}
	})
	app.Lifecycle().SetOnStopped(func() {
		fmt.Printf("%+v\n", mainWindow.Canvas().Size())
	})
	mainWindow.Show()
	app.Run()
}
