package core

import (
	"fmt"
	"net/http"
)

func (c *Client) setHeader(http *http.Request) {

	userAgent, os := randomUserAgent()
	if userAgent == "" || os == "" {
		userAgent = "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.2 Chrome/38.0.2125.102 Mobile Safari/537.36"
		os = "Android"
	}

	header := map[string]string{
		"accept":             "*/*",
		"content-type":       "application/json",
		"accept-language":    "en-US,en;q=0.9",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Android WebView";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": fmt.Sprintf("\"%s\"", os),
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"Referer":            "https://cf.seeddao.org/",
		"Referrer-Policy":    "strict-origin-when-cross-origin",
		"X-Requested-With":   "org.telegram.messenger.web",
		"User-Agent":         userAgent,
	}

	if c.AccessToken != "" {
		header["telegram-data"] = c.AccessToken
	}

	for key, value := range header {
		http.Header.Set(key, value)
	}
}
