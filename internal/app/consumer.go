package app

import (
	"fmt"
	"fourth-exam/user-service-evrone/internal/delivery/kafka/handlers"
	"fourth-exam/user-service-evrone/internal/infrastructure/kafka"
	"fourth-exam/user-service-evrone/internal/infrastructure/repository/postgresql"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"fourth-exam/user-service-evrone/internal/pkg/postgres"
	"fourth-exam/user-service-evrone/internal/usecase"
	"fourth-exam/user-service-evrone/internal/usecase/event"
	"time"

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
	fmt.Println(err, " error here")
	if err != nil {
		fmt.Println(err, " error here")
		return nil, err
	}

	return &UserConsumer{Config: conf, Logger: logger, DB: db, BrokerConsumer: consumer}, nil
}

func (u *UserConsumer) Run() error {

	// repo init
	userRepo := postgresql.NewUsersRepo(u.DB)

	// usecase init
	duration, err := time.ParseDuration(u.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	userUseCase := usecase.NewUserService(duration, userRepo)

	// event handler
	eventHandler := handlers.NewUserConsumerHandler(u.Config, u.BrokerConsumer, u.Logger, userUseCase)

	return eventHandler.HandlerEvents()
}

func (u *UserConsumer) Close() {
	u.BrokerConsumer.Close()

	u.Logger.Sync()
}
