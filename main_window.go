package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	tabTitlePages        = "Page"
	tabTitleNotification = "Notification"
)

func setupMainWindow(app fyne.App) fyne.Window {
	mainWindow := app.NewWindow(app.Metadata().Name)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(tabTitlePages, theme.DocumentIcon(), pagesTab()),
		container.NewTabItemWithIcon(tabTitleNotification, theme.MailSendIcon(), notificationTab()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	mainWindow.SetContent(tabs)
	return mainWindow
}

func pagesTab() fyne.CanvasObject {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://example.com/")
	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder("1024")
	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("1280")

	form := widget.NewForm(
		widget.NewFormItem("URL", urlEntry),
		widget.NewFormItem("Width", widthEntry),
		widget.NewFormItem("Height", heightEntry),
	)
	return container.NewVBox(
		widget.NewLabelWithStyle(tabTitlePages, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
	)
}

func notificationTab() fyne.CanvasObject {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://example.com/")

	form := widget.NewForm(
		widget.NewFormItem("URL", urlEntry),
	)
	return container.NewVBox(
		widget.NewLabelWithStyle(tabTitleNotification, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		form,
	)
}
