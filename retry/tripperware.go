// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package http_retry

import (
	"net/http"
	"time"

	"context"
	"errors"
	"fmt"

	"github.com/mwitkow/go-httpwares"
)

// Tripperware is client side HTTP ware that retries the requests.
//
// Be default this retries safe and idempotent requests 3 times with a linear delay of 100ms. This behaviour can be
// customized using With* parameter options.
//
// Requests that have `http_retry.Enable` set on them will always be retried.
func Tripperware(opts ...Option) httpwares.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		o := evaluateOptions(opts)
		return httpwares.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			// Short-circuit to avoid allocations.
			if !o.decider(req) && !isEnabled(req.Context()) {
				return next.RoundTrip(req)
			}
			if o.maxRetry == 0 || req.GetBody == nil {
				// If we are configured to do no retries or the lack of GetBody function doesn't allow for re-reads of
				// body data.
				return next.RoundTrip(req)
			}
			var err error
			var lastResp *http.Response
			for attempt := uint(0); attempt < o.maxRetry; attempt++ {
				thisReq := req.WithContext(req.Context()) // make a copy.
				thisReq.Body, err = req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("failed reading body for retry: %v", err)
				}
				if err := waitRetryBackoff(attempt, req.Context(), o); err != nil {
					return nil, err // context errors from req.Context()
				}
				lastResp, err = next.RoundTrip(thisReq)
				if err == nil {
					if !o.discarder(lastResp) {
						return lastResp, nil
					}
				} else if isContextError(err) {
					return nil, err
				}
			}
			if lastResp != nil {
				return lastResp, err
			} else if err != nil {
				return nil, err
			}
			return nil, errors.New("maximum retry budget reached")
		})
	}
}

func waitRetryBackoff(attempt uint, parentCtx context.Context, opt *options) error {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = opt.backoffFunc(attempt)
	}
	if waitTime > 0 {
		select {
		case <-parentCtx.Done():
			return parentCtx.Err()
		case <-time.Tick(waitTime):
		}
	}
	return nil
}

func isContextError(err error) bool {
	return err == context.DeadlineExceeded || err == context.Canceled
}