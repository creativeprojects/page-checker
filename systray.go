package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func setupSystrayMenu(app fyne.App, window fyne.Window) {
	if desk, ok := app.(desktop.App); ok {
		m := fyne.NewMenu(app.Metadata().Name,
			fyne.NewMenuItem("Show", func() {
				window.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}
}
