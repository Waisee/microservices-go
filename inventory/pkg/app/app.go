package app

import (
	"log/slog"

	"google.golang.org/grpc"

	inventoryapi "github.com/waisee/microservices-go/inventory/internal/api/inventory/v1"
	interceptor "github.com/waisee/microservices-go/inventory/internal/interceptor"
	partrepo "github.com/waisee/microservices-go/inventory/internal/repository/part"
	partsvc "github.com/waisee/microservices-go/inventory/internal/service/part"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptor.UnaryServerInterceptor(slog.Default())),
	}
}

func RegisterServices(grpcServer *grpc.Server) {
	repo := partrepo.NewPartRepository()
	svc := partsvc.NewPartService(repo)
	api := inventoryapi.NewInventoryApi(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
