package main

import (
	"net"
	"net/url"
	"net/http"
	"time"
	"errors"
	"io"
)

type UrlProcessor struct {
	client *http.Client
}

func NewUrlProcessor(timeout int) *UrlProcessor {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return &UrlProcessor {
		client: client,
	}
}

func (u *UrlProcessor) ProcessUrl(urlArg string) Result {
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
	resp, err := u.client.Get(urlString)

	if err != nil {
		var dnsErr *net.DNSError
		var netErr net.Error
		if errors.As(err, &dnsErr) {
			result.Error = "DNS lookup failed"
		} else if errors.As(err, &netErr) && netErr.Timeout() {
			result.Error = "Request timeout"
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
