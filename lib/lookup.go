package lib

import (
	"context"
	"fmt"
	"net"
	"net/url"
)

func LookupHostname(ctx context.Context, cfg Flags) error {
	sourceURL, err := url.Parse(cfg.URL)
	if err != nil {
		return NewError("cannot parse URL", ErrorCodeNameNotResolved, err)
	}
	hostname := sourceURL.Hostname()
	resolver := net.DefaultResolver
	cname, err := resolver.LookupCNAME(ctx, hostname)
	if err != nil {
		return NewError("cannot resolve hostname", ErrorCodeNameNotResolved, err)
	}
	if cfg.Verbose && len(cname) > 0 {
		fmt.Printf("Host %q resolving to %q\n", hostname, cname)
	}
	_, err = resolver.LookupHost(ctx, hostname)
	if err != nil {
		return NewError("cannot resolve hostname", ErrorCodeNameResolutionFailed, err)
	}
	return nil
}
