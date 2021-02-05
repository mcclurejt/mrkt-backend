package iex

import (
	"time"

	"github.com/sirupsen/logrus"
)

type ClientOption func(*RealClient) *RealClient

func ClientOptionSetTimeout(timeout time.Duration) ClientOption {
	return func(r *RealClient) *RealClient {
		r.Log.WithFields(logrus.Fields{"timeout": timeout}).Info("Set client timeout")
		r.HTTPClient.Timeout = timeout
		return r
	}
}

func ClientOptionSetLogLevel(level logrus.Level) ClientOption {
	return func(r *RealClient) *RealClient {
		r.Log.WithFields(logrus.Fields{"logLevel": level.String()}).Info("Set log level")
		r.Log.SetLevel(level)
		return r
	}
}
