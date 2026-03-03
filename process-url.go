package main

import (
	"net"
	"net/url"
	"net/http"
	"time"
	"errors"
	"io"
	"fmt"
	"context"
)

type UrlProcessor struct {
	timeout int
	client  *http.Client
}

func NewUrlProcessor(timeout int) *UrlProcessor {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return &UrlProcessor {
		client : client,
		timeout: timeout,
	}
}

func (u *UrlProcessor) ProcessUrl(ctx context.Context, urlArg string) Result {
	var result Result

	reqUrl, err := url.ParseRequestURI(urlArg)
	if err != nil {
		fmt.Println(urlArg)
		result.URL = "bad-url"
		result.Error = "invalid URL"
		return result
	}

	urlString := reqUrl.String()

	result.URL = urlString

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		fmt.Println("error forming request", err.Error())
	}

	startTime := time.Now()
	// resp, err := u.client.Get(urlString)
	resp, err := u.client.Do(req)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = "cancelled (global timeout)"
		} else if ctx.Err() == context.Canceled {
			result.Error = "cancelled (user interrupted)"
		} else {
			var dnsErr *net.DNSError
			var netErr net.Error
			if errors.As(err, &dnsErr) {
				result.Error = "DNS lookup failed"
			} else if errors.As(err, &netErr) && netErr.Timeout() {
				result.Error = fmt.Sprintf("request timeout after %ds", int(time.Since(startTime).Seconds()))
			} else {
				result.Error = err.Error()
			}
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
