package usecase

import (
	"context"
	"time"

	"fourth-exam/user-service-evrone/internal/entity"
	"fourth-exam/user-service-evrone/internal/infrastructure/repository"
)

type User interface {
	Create(ctx context.Context, req *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, req *entity.GetListFilter) ([]*entity.User, error)
	Update(ctx context.Context, req *entity.User) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	BaseUseCase
	repo       repository.User
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.User) User {
	return &userService{
		repo:       repo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *userService) Create(ctx context.Context, req *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

	return u.repo.Create(ctx, req)
}

func (u *userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Get(ctx, params)
}

func (u *userService) List(ctx context.Context, req *entity.GetListFilter) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.List(ctx, req)
}

func (u *userService) Update(ctx context.Context, req *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&req.Id, &req.CreatedAt, &req.UpdatedAt)

	return u.repo.Update(ctx, req)
}

func (u *userService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
