package repository
import (
	"context"
	"fourth-exam/user-service-evrone/internal/entity"
)

type User interface {
	Create(ctx context.Context, req *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, req *entity.GetListFilter) ([]*entity.User, error)
	Update(ctx context.Context, req *entity.User) (error)
	Delete(ctx context.Context, id string) error
}