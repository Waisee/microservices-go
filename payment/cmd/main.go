package main

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	app "github.com/waisee/microservices-go/payment/pkg/app"
)

const grpcAddress = "127.0.0.1:50052"
const (
	maxConnectionIdle     = 15 * time.Minute
	maxConnectionAge      = 30 * time.Minute
	maxConnectionAgeGrace = 5 * time.Second
	keepaliveTime         = 5 * time.Minute
	timeout               = 1 * time.Second
	minTime               = 5 * time.Minute
)

func main() {
	//nolint:noctx // Контекст здесь не нужен: GracefulStop() сам закроет listener и прервёт Accept()
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("ошибка запуска слушателя", "error", err)
		return
	}
	// Примечание: defer lis.Close() не нужен, так как GracefulStop() закрывает listener автоматически

	grpcServer := grpc.NewServer(
		append(app.Interceptors(),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     maxConnectionIdle,
				MaxConnectionAge:      maxConnectionAge,
				MaxConnectionAgeGrace: maxConnectionAgeGrace,
				Time:                  keepaliveTime,
				Timeout:               timeout,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             minTime,
				PermitWithoutStream: true,
			}),
		)...,
	)

	app.RegisterServices(grpcServer)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	slog.Info("запуск PaymentService", "адрес", grpcAddress)

	// Контекст, который отменяется по SIGINT/SIGTERM или при падении сервера.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("🚀 gRPC сервер запущен", "address", grpcAddress)
		err = grpcServer.Serve(lis)
		if err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
			cancel()
		}
	}()

	// Ждём сигнал от ОС или падение сервера.
	<-ctx.Done()
	slog.Info("🛑 остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("✅ сервер остановлен")
}
