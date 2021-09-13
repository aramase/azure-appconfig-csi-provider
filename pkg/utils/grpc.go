package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// ParseEndpoint parses the endpoint string and returns the protocol, endpoint and error.
func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid endpoint: %v", ep)
}

// LogInterceptor is a gRPC interceptor that logs the gRPC requests and responses.
func LogInterceptor(log logr.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		ctxDeadline, _ := ctx.Deadline()
		log.Info("request", "method", info.FullMethod, "deadline", time.Until(ctxDeadline).String())

		resp, err := handler(ctx, req)
		s, _ := status.FromError(err)
		log.Info("response", "method", info.FullMethod, "duration", time.Since(start).String(), "code", s.Code().String(), "message", s.Message())

		return resp, err
	}
}
