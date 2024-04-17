package app

import (
	"fmt"
	pb "fourth-exam/user-service-evrone/genproto/user_service"
	"fourth-exam/user-service-evrone/internal/delivery/grpc/server"
	"fourth-exam/user-service-evrone/internal/delivery/grpc/services"
	grpc_service_clients "fourth-exam/user-service-evrone/internal/infrastructure/grpc_service_client"
	repo "fourth-exam/user-service-evrone/internal/infrastructure/repository/postgresql"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"fourth-exam/user-service-evrone/internal/pkg/logger"
	"fourth-exam/user-service-evrone/internal/pkg/postgres"
	"fourth-exam/user-service-evrone/internal/usecase"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	ServiceClients grpc_service_clients.ServiceClients
	GrpcServer     *grpc.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	clients, err := grpc_service_clients.New(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		Logger: logger,
		DB:     db,
		GrpcServer: grpcServer,
		ServiceClients: clients,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}

	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	userRepo := repo.NewUsersRepo(a.DB)

	userUseCase := usecase.NewUserService(contextTimeout, userRepo)

	pb.RegisterUserServiceServer(a.GrpcServer, services.NewRPC(a.Logger, userUseCase))

	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// zap logger sync
	a.Logger.Sync()
}
