package app

import (
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
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

func RegisterServices(grpcServer *grpc.Server, pool *pgxpool.Pool) {
	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("не удалось создать менеджер транзакций", "error", err)
		os.Exit(1)
	}
	repo := partrepo.NewPartRepository(pool, txManager)
	svc := partsvc.NewPartService(repo)
	api := inventoryapi.NewInventoryApi(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
