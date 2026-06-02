package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
)

func OgenErrorHandler(logger *slog.Logger) ogenerrors.ErrorHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			logger.ErrorContext(ctx, "ogen handler error", "error", err)
		}
		ogenerrors.DefaultErrorHandler(ctx, w, r, err)
	}
}
