package services

import (
	"context"
	pb "fourth-exam/user-service-evrone/genproto/user_service"
	grpc "fourth-exam/user-service-evrone/internal/delivery"
	"fourth-exam/user-service-evrone/internal/entity"
	"fourth-exam/user-service-evrone/internal/usecase"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userRPC struct {
	logger      *zap.Logger
	userUsecase usecase.User
}

func NewRPC(logger *zap.Logger, userUsecase usecase.User) pb.UserServiceServer {
	return &userRPC{
		logger:      logger,
		userUsecase: userUsecase,
	}
}

func (d *userRPC) Create(ctx context.Context, in *pb.User) (*pb.User, error) {
	id := uuid.New().String()
	_, err := d.userUsecase.Create(ctx, &entity.User{
		Id:           id,
		Email:        in.Email,
		Password:     in.Password,
		Username:     in.Username,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Bio:          in.Bio,
		Website:      in.Bio,
		RefreshToken: in.RefreshToken,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		d.logger.Error("userUseCase.Create", zap.Error(err))
		return &pb.User{}, grpc.Error(ctx, err)
	}
	in.Id = id
	return in, nil
}

func (d *userRPC) Update(ctx context.Context, in *pb.User) (*pb.User, error) {
	err := d.userUsecase.Update(ctx, &entity.User{
		Id:           in.Id,
		Email:        in.Email,
		Password:     in.Password,
		Username:     in.Username,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Bio:          in.Bio,
		Website:      in.Bio,
		RefreshToken: in.RefreshToken,
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		d.logger.Error("userUseCase.Update", zap.Error(err))
		return &pb.User{}, grpc.Error(ctx, err)
	}

	return in, nil
}

func (d *userRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.UserModel, error) {
	user, err := d.userUsecase.Get(ctx, map[string]string{"id": in.UserId})
	if err != nil {
		d.logger.Error("userUseCase.Get", zap.Error(err))
		return &pb.UserModel{}, grpc.Error(ctx, err)
	}

	return &pb.UserModel{
		Id:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bio:       user.Bio,
		Website:   user.Website,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (d *userRPC) Delete(ctx context.Context, in *pb.GetRequest) (*empty.Empty, error) {
	err := d.userUsecase.Delete(ctx, in.UserId)
	if err != nil {
		d.logger.Error("userUseCase.Delete", zap.Error(err))
		return &empty.Empty{}, grpc.Error(ctx, err)
	}

	return &empty.Empty{}, nil
}

func (d *userRPC) List(ctx context.Context, in *pb.GetListFilter) (*pb.Users, error) {
	filter := &entity.GetListFilter{
		Limit:   in.Limit,
		Page:    in.Page,
		OrderBy: in.OrderBy,
	}

	users, err := d.userUsecase.List(ctx, filter)
	if err != nil {
		d.logger.Error("userUseCase.List", zap.Error(err))
		return &pb.Users{}, grpc.Error(ctx, err)
	}

	var pbUsers []*pb.UserModel
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.UserModel{
			Id:        user.Id,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Bio:       user.Bio,
			Website:   user.Website,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
			Posts:     []*pb.Post{},
		})
	}

	return &pb.Users{Users: pbUsers}, nil
}
