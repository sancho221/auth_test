package handler

import (
	"auth_test/pkg/metrics"
	"context"
	"time"

	"google.golang.org/grpc"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		method := extractMethodName(info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		status := "success"
		if err != nil {
			status = "error"
		}

		metrics.GRPCRequests.WithLabelValues(method, status).Inc()
		metrics.GRPCRequestDuration.WithLabelValues(method).Observe(duration)

		return resp, err
	}
}

func extractMethodName(fullMethod string) string {
	for i := len(fullMethod) - 1; i >= 0; i-- {
		if fullMethod[i] == '/' {
			return fullMethod[i+1:]
		}
	}
	return fullMethod
}
