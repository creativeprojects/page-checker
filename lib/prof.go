package lib

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func StartProfiler(cfg Flags) func() error {
	noClose := func() error { return nil }
	if cfg.ProfilerPort == 0 {
		return noClose
	}
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	server := &http.Server{
		Addr: fmt.Sprintf("localhost:%d", cfg.ProfilerPort),
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	return server.Close
}
