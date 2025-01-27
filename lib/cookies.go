package lib

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-rod/rod/lib/proto"
)

func ConvertCookies(source []string) (cookies []*proto.NetworkCookieParam) {
	for _, sourceCookie := range source {
		decoder := json.NewDecoder(strings.NewReader(sourceCookie))
		cookie := &proto.NetworkCookieParam{}
		err := decoder.Decode(cookie)
		if err != nil {
			fmt.Printf("error decoding cookie: %s\n", err)
			continue
		}
		if cookie.Name == "" {
			continue
		}
		if cookie.Path == "" {
			cookie.Path = "/"
		}
		cookies = append(cookies, cookie)
	}
	return cookies
}
