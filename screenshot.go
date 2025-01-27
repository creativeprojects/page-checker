package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/creativeprojects/page-checker/lib"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func screenshot() error {
	var err error
	var cfg lib.Flags

	ctx := context.Background()
	{
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		err = lib.LookupHostname(ctx, cfg)
		if err != nil {
			return err
		}
	}

	browser, err := lib.NewBrowser(cfg)
	if err != nil {
		return err
	}
	defer browser.Close()

	page, err := lib.NewPage(cfg, browser)
	if err != nil {
		return err
	}
	defer page.Close()

	result := make(chan lib.Response, 1)
	go func() {
		lib.PageResponse(cfg, page, result)
	}()

	waitNavigation := pageReady(cfg, page)

	err = page.Navigate(cfg.URL)
	if err != nil {
		code := lib.GetCodeFromError(err)
		if code == 0 {
			code = lib.ErrorCodeAddressUnreachable
		}
		return lib.NewError("cannot load webpage", code, err)
	}

	// keep a maximum timeout of 10 minutes
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 600
	}
	ctx = context.Background()
	{
		ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
		err = lib.RunWithContext(ctx, waitNavigation)
	}
	if err != nil {
		// we keep going if savepartial is set
		if !(cfg.Savepartial && errors.Is(err, context.DeadlineExceeded)) {
			code := lib.GetCodeFromError(err)
			if code == 0 {
				code = lib.ErrorCodeUnknown
			}
			return lib.NewError("error waiting for the page to finish loading", code, err)
		}
	}

	if len(cfg.ClickXPath) > 0 {
		for _, xpath := range cfg.ClickXPath {
			clickXPath(page, xpath)
		}
	}

	time.Sleep(time.Duration(cfg.Sleep) * time.Second)

	<-result

	req := proto.PageCaptureScreenshot{
		Format:      proto.PageCaptureScreenshotFormatPng,
		Quality:     &cfg.Quality,
		FromSurface: false,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  float64(cfg.Width),
			Height: float64(cfg.Height),
			Scale:  cfg.Scale,
		},
	}
	content, err := page.Screenshot(true, &req)
	if err != nil {
		code := lib.GetCodeFromError(err)
		if code == 0 {
			code = lib.ErrorCodeEmptyResponse
		}
		return lib.NewError("cannot render a screenshot", code, err)
	}
	if cfg.Output != "" {
		err = os.WriteFile(cfg.Output, content, 0666)
		if err != nil {
			return lib.NewError("cannot save file", lib.ErrorCodeCannotSaveFile, err)
		}
	}
	return nil
}

func clickXPath(page *rod.Page, xpath string) {
	element, err := page.ElementX(xpath)
	if err != nil {
		fmt.Printf("element %q: %s\n", xpath, err)
		return
	}
	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		fmt.Printf("click on element: %s\n", err)
		return
	}
}
