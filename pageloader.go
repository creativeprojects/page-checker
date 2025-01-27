package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/creativeprojects/page-checker/lib"
	"github.com/go-rod/rod/lib/proto"
)

func pageloader() error {
	var err error
	var cfg lib.Flags

	closeProfiler := lib.StartProfiler(cfg)
	defer closeProfiler()

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

	// adds 1 minute for the browser timeout
	timeout := time.Duration(cfg.Timeout+60) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	browser = browser.Context(ctx)

	page, err := lib.NewPage(cfg, browser)
	if err != nil {
		return err
	}
	defer page.Close()

	result := make(chan lib.Response, 1)
	go func() {
		lib.PageResponse(cfg, page, result)
	}()
	// sends an empty response if not replied in time
	time.AfterFunc(time.Duration(cfg.Timeout)*time.Second, func() {
		result <- lib.Response{}
	})

	waitNavigation := page.WaitNavigation(proto.PageLifecycleEventNameDOMContentLoaded)

	err = page.Navigate(cfg.URL)
	if err != nil {
		code := lib.GetCodeFromError(err)
		if code == 0 {
			code = lib.ErrorCodeAddressUnreachable
		}
		return lib.NewError("cannot load webpage", code, err)
	}

	waitNavigation()

	// we keep going if savepartial is set
	if !cfg.Savepartial && err != nil {
		code := lib.GetCodeFromError(err)
		if code == 0 {
			code = lib.ErrorCodeUnknown
		}
		return lib.NewError("error waiting for webpage to finish loading", code, err)
	}
	time.Sleep(time.Duration(cfg.Sleep) * time.Second)

	content, err := page.HTML()
	if err != nil {
		code := lib.GetCodeFromError(err)
		if code == 0 {
			code = lib.ErrorCodeEmptyResponse
		}
		return lib.NewError("cannot load HTML content", code, err)
	}
	if cfg.Output != "" {
		err = os.WriteFile(cfg.Output, []byte(content), 0666)
		if err != nil {
			return lib.NewError("cannot save file", lib.ErrorCodeCannotSaveFile, err)
		}

		response := <-result
		encoder := json.NewEncoder(os.Stdout)
		err = encoder.Encode(&response)
		if err != nil {
			return lib.NewError("cannot generate JSON result", lib.ErrorCodeOther, err)
		}
	} else {
		fmt.Println(content)
	}
	return nil
}
