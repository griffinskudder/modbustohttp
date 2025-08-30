package interceptors

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
)

func NewLoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
			logger.Info(
				"Calling procedure",
				slog.String("procedure", request.Spec().Procedure),
				slog.String("request", fmt.Sprintf("%v", request.Any())),
			)
			response, err := next(ctx, request)
			if err != nil {
				logger.Error(
					"Error calling procedure",
					slog.String("procedure", request.Spec().Procedure),
					slog.String("request", fmt.Sprintf("%v", request.Any())),
					slog.String("error", err.Error()),
				)
			} else {
				logger.Info(
					"Procedure call successful",
					slog.String("procedure", request.Spec().Procedure),
					slog.String("request", fmt.Sprintf("%v", request.Any())),
					slog.String("response", fmt.Sprintf("%v", response.Any())),
				)
			}
			return response, err
		}
	}
}
