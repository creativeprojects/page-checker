package main

import (
	"github.com/creativeprojects/page-checker/lib"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func pageReady(cfg lib.Flags, page *rod.Page) func() error {
	switch {
	case cfg.Nowait:
		return func() error {
			return page.WaitLoad()
		}
	case cfg.Noidlewait:
		_ = proto.PageSetLifecycleEventsEnabled{Enabled: true}.Call(page)

		wait1 := page.EachEvent(func(e *proto.PageLifecycleEvent) bool {
			return e.Name == proto.PageLifecycleEventNameDOMContentLoaded
		})

		wait2 := page.EachEvent(func(e *proto.PageLifecycleEvent) bool {
			return e.Name == proto.PageLifecycleEventNameNetworkAlmostIdle
		})

		return func() error {
			wait1()
			wait2()
			_ = proto.PageSetLifecycleEventsEnabled{Enabled: false}.Call(page)
			return nil
		}
	default:
		wait := page.WaitNavigation(proto.PageLifecycleEventNameNetworkIdle)
		return func() error {
			err := page.WaitLoad()
			if err != nil {
				return err
			}
			wait()
			return nil
		}
	}
}
