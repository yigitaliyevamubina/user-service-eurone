package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"fourth-exam/user-service-evrone/internal/entity"
	"fourth-exam/user-service-evrone/internal/pkg/otlp"
	"fourth-exam/user-service-evrone/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	usersTableName     = "users"
	userServiceName    = "userService"
	userSpanRepoPrefix = "userServiceRepo"
)

type userRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) *userRepo {
	return &userRepo{
		tableName: usersTableName,
		db:        db,
	}
}

func (u *userRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"username",
		"email",
		"password",
		"first_name",
		"last_name",
		"bio",
		"website",
		"is_active",
		"refresh_token",
		"created_at",
		"updated_at",
	).From(u.tableName)
}

func (u *userRepo) Create(ctx context.Context, req *entity.User) (*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Create")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "User -> repository -> ", Value: attribute.StringValue("Create user")})

	data := map[string]any{
		"id":            req.Id,
		"username":      req.Username,
		"email":         req.Email,
		"password":      req.Password,
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"bio":           req.Bio,
		"website":       req.Website,
		"is_active":     req.IsActive,
		"refresh_token": req.RefreshToken,
		"created_at":    req.CreatedAt,
		"updated_at":    req.UpdatedAt,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, "create"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
}

func (u *userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Get")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "User -> repository -> ", Value: attribute.StringValue("Get user")})

	var (
		user entity.User
	)

	queryBuilder := u.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, "get"))
	}

	var (
		updatedAt sql.NullTime
	)
	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.Website,
		&user.IsActive,
		&user.RefreshToken,
		&user.CreatedAt,
		&updatedAt,
	); err != nil {
		return nil, u.db.Error(err)
	}

	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	return &user, nil
}

func (u *userRepo) List(ctx context.Context, req *entity.GetListFilter) ([]*entity.User, error) {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"List")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "User -> repository -> ", Value: attribute.StringValue("Get list")})

	var (
		users []*entity.User
	)
	queryBuilder := u.usersSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	if req.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(req.OrderBy)
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, "list"))
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()
	var (
		updatedAt sql.NullTime
	)
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Bio,
			&user.Website,
			&user.IsActive,
			&user.RefreshToken,
			&user.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, u.db.Error(err)
		}

		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *userRepo) Update(ctx context.Context, req *entity.User) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Update")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "User -> repository -> ", Value: attribute.StringValue("Update user")})

	data := map[string]any{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"username":   req.Username,
		"bio":        req.Bio,
		"website":    req.Website,
		"is_active":  req.IsActive,
		"updated_at": req.UpdatedAt,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.tableName).
		SetMap(data).
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.tableName+" update")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (u *userRepo) Delete(ctx context.Context, id string) error {
	ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Delete")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "User -> repository -> ", Value: attribute.StringValue("Delete user")})

	sqlStr, args, err := u.db.Sq.Builder.
		Delete(u.tableName).
		Where(u.db.Sq.Equal("id", id)).
		ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.tableName+" delete")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}
