package app

import (
	"log/slog"

	"google.golang.org/grpc"

	v1 "github.com/waisee/microservices-go/payment/internal/api/payment/v1"
	interceptor "github.com/waisee/microservices-go/payment/internal/interceptor"
	payment "github.com/waisee/microservices-go/payment/internal/service/payment"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptor.UnaryServerInterceptor(slog.Default())),
	}
}

func RegisterServices(grpcServer *grpc.Server) {
	paymentService := payment.NewPaymentService()
	paymentAPI := v1.NewPaymentAPI(paymentService)
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentAPI)
}
