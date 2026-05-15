package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waisee/microservices-go/inventory/internal/errors"
)

func UnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = slog.Default()
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}

		switch {
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, errs.ErrPartNotFound.Error())
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Error(codes.InvalidArgument, errs.ErrInvalidUUID.Error())
		}

		logger.Error("unexpected grpc error", "method", info.FullMethod, "error", err)
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
