package app

import (
	"fourth-exam/user-service-evrone/internal/delivery/kafka/handlers"
	"fourth-exam/user-service-evrone/internal/infrastructure/kafka"
	"fourth-exam/user-service-evrone/internal/infrastructure/repository/postgresql"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"fourth-exam/user-service-evrone/internal/pkg/postgres"
	"fourth-exam/user-service-evrone/internal/usecase"
	"fourth-exam/user-service-evrone/internal/usecase/event"

	logpkg "fourth-exam/user-service-evrone/internal/pkg/logger"

	"go.uber.org/zap"
)

type UserConsumer struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewUserConsumer(conf *config.Config) (*UserConsumer, error) {
	logger, err := logpkg.New(conf.LogLevel, conf.Environment, conf.APP+"_cli"+".lo")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(conf)
	if err != nil {
		return nil, err
	}

	return &UserConsumer{Config: conf, Logger: logger, DB: db, BrokerConsumer: consumer}, nil
}

func (u *UserConsumer) Run() error {
	
	// repo init
	userRepo := postgresql.NewUsersRepo(u.DB)

	// usecase init
	userUseCase := usecase.NewUserService(u.DB.Config().ConnConfig.ConnectTimeout, userRepo)

	// event handler
	eventHandler := handlers.NewUserConsumerHandler(u.Config, u.BrokerConsumer, u.Logger, userUseCase)

	return eventHandler.HandlerEvents()
}

func (u *UserConsumer) Close() {
	u.BrokerConsumer.Close()

	u.Logger.Sync()
}