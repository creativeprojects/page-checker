package lib

import (
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Response struct {
	Status   int               `json:"status"`
	Headers  map[string]string `json:"headers"`
	Redirect string            `json:"redirect"`
}

func PageResponse(cfg Flags, page *rod.Page, sendResponse chan Response) {
	wait := page.EachEvent(
		func(e *proto.NetworkResponseReceived) bool {
			if strings.HasPrefix(e.Response.URL, "data:") {
				return false
			}
			Verbose(cfg, "## network response: id=%q url=%q statusCode=%d \n", e.RequestID, e.Response.URL, e.Response.Status)
			if e.Type != proto.NetworkResourceTypeDocument {
				return false
			}
			response := Response{
				Status:   e.Response.Status,
				Headers:  map[string]string{},
				Redirect: e.Response.URL,
			}

			for key, value := range e.Response.Headers {
				response.Headers[key] = value.String()
			}
			sendResponse <- response
			return true
		},
	)
	wait()
}
