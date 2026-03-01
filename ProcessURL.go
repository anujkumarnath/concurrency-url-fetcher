package main

import (
	"net"
	"net/url"
	"net/http"
	"time"
	"errors"
	"io"
)

func ProcessUrl(urlArg string) Result {
	var result Result

	reqUrl, err := url.ParseRequestURI(urlArg)
	if err != nil {
		result.URL = "bad-url"
		result.Error = "invalid URL"
		return result
	}

	urlString := reqUrl.String()

	result.URL = urlString

	startTime := time.Now()
	resp, err := http.Get(urlString)

	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			result.Error = "DNS lookup failed"
		} else {
			result.Error = err.Error()
		}
		return result
	}

	defer resp.Body.Close()

	endTime  := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		result.Duration = duration
		return result
	}

	result.Status   = resp.Status
	result.Size     = len(body)
	result.Duration = duration

	return result
}
