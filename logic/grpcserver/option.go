package grpcserver

import (
	"math"
	"path"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
)

func newDefaultServerOptions() *options {
	return &options{
		maxConcurrentStreams: math.MaxUint32,
		loggingDecider: func(fullMethodName string, err error) bool {
			service := path.Dir(fullMethodName)[1:]
			return service != "grpc.reflection.v1alpha.ServerReflection"
		},
	}
}

type Option func(*options)

type options struct {
	maxConcurrentStreams uint32
	loggingDecider       grpcLogging.Decider
}

func WithMaxConcurrentStreams(m uint32) Option {
	return func(o *options) {
		o.maxConcurrentStreams = m
	}
}

func WithLoggingDecider(d grpcLogging.Decider) Option {
	return func(o *options) {
		o.loggingDecider = d
	}
}
