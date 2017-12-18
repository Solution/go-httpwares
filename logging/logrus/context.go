// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package http_logrus

import (
	"net/http"

	"github.com/improbable-eng/go-httpwares/tags/logrus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"io/ioutil"
)

var (
	nullLogger = &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.PanicLevel,
	}
)

// Extract takes the call-scoped logrus.Entry from grpc_logrus middleware.
//
// The logger will have fields pre-populated using http_ctxtags.
//
// If the http_logrus middleware wasn't used, a no-op `logrus.Entry` is returned. This makes it safe to use regardless.
// Deprecated use ctx_logrus.Extract instead
func Extract(req *http.Request) *logrus.Entry {
	return ctx_logrus.ExtractFromContext(req.Context())
}

// Extract takes the call-scoped logrus.Entry from grpc_logrus middleware.
//
// The logger will have fields pre-populated using http_ctxtags.
//
// If the http_logrus middleware wasn't used, a no-op `logrus.Entry` is returned. This makes it safe to use regardless.
// Deprecated use ctx_logrus.ExtractFromContext instead
func ExtractFromContext(ctx context.Context) *logrus.Entry {
	return ctx_logrus.ExtractFromContext(ctx)
}
