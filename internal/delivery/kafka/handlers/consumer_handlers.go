package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	pb "fourth-exam/user-service-evrone/genproto/user_service"
	"fourth-exam/user-service-evrone/internal/entity"
	"fourth-exam/user-service-evrone/internal/infrastructure/kafka"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"fourth-exam/user-service-evrone/internal/usecase"
	"fourth-exam/user-service-evrone/internal/usecase/event"
	"time"

	"go.uber.org/zap"
	"github.com/k0kubun/pp"
)

type userConsumerHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	userUsecase    usecase.User
}

func NewUserConsumerHandler(conf *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	userUseCase usecase.User) *userConsumerHandler {
	return &userConsumerHandler{
		config:         conf,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		userUsecase:    userUseCase,
	}
}

func (u *userConsumerHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		u.config.Kafka.Address,
		u.config.Kafka.Topic.UserTopic,
		"1",
		func(ctx context.Context, key, value []byte) error {
			var user pb.User

			if err := json.Unmarshal(value, &user); err != nil {
				return err
			}

			pp.Println(user)

			req := entity.User{
				Id:           user.Id,
				Username:     user.Username,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Bio:          user.Bio,
				Website:      user.Website,
				IsActive:     user.IsActive,
				RefreshToken: user.RefreshToken,
				Email:        user.Email,
			}

			ctxNew, err := context.WithTimeout(context.Background(), time.Second*7)
			if err != nil {
				fmt.Println(err, "Context error")
			}
			_, errr := u.userUsecase.Create(ctxNew, &req)
			if errr != nil {
				fmt.Println(errr, "Create=========================")
			}
			// fmt.Println(req, "user")

			return nil
		},
	)

	u.brokerConsumer.RegisterConsumer(consumerConfig)
	u.brokerConsumer.Run()

	return nil
}
