package lib

import (
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func NewBrowser(cfg Flags) (*rod.Browser, error) {
	var err error

	browser := rod.New()
	if !cfg.Verbose {
		browser.Logger(&Logger{})
	}
	err = browser.Connect()
	if err != nil {
		return nil, NewError("cannot launch chrome", ErrorCodeCannotLaunchChrome, err)
	}
	return browser, nil
}

func NewPage(cfg Flags, browser *rod.Browser) (*rod.Page, error) {
	var err error

	version, err := browser.Version()
	if err != nil {
		return nil, NewError("cannot detect browser version", ErrorCodeCannotLaunchChrome, err)
	}
	Verbose(cfg, "User-Agent: %s\n", version.UserAgent)

	userAgent := cfg.Useragent
	if userAgent == "" {
		userAgent = strings.ReplaceAll(version.UserAgent, "Headless", "")
	}
	Verbose(cfg, "Using:      %s\n", userAgent)

	var page *rod.Page
	if cfg.Stealth {
		page, err = stealth.Page(browser)
		if err != nil {
			return nil, NewError("cannot create a page in stealth mode", ErrorCodeCannotLaunchChrome, err)
		}
	} else {
		page, err = browser.Page(proto.TargetCreateTarget{})
		if err != nil {
			return nil, NewError("cannot create a page", ErrorCodeCannotLaunchChrome, err)
		}
	}
	if cfg.Timeout > 0 {
		// add 30 seconds to allow for the savepartial flag (page.HTML() would get cancelled otherwise)
		page = page.Timeout(time.Duration(cfg.Timeout+30) * time.Second)
	}

	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, NewError("cannot set user agent", ErrorCodeCannotLaunchChrome, err)
	}
	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             cfg.Width,
		Height:            cfg.Height,
		DeviceScaleFactor: 2,
	})
	if err != nil {
		return nil, NewError("cannot set viewport", ErrorCodeCannotLaunchChrome, err)
	}
	err = page.SetCookies(ConvertCookies(cfg.Cookie))
	if err != nil {
		return nil, NewError("cannot set cookies", ErrorCodeCannotLaunchChrome, err)
	}
	return page, nil
}
