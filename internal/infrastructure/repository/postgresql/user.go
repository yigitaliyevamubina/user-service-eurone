package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"fourth-exam/user-service-evrone/internal/entity"
	"fourth-exam/user-service-evrone/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName  = "users"
	userServiceName = "userService"
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
	fmt.Println("2")
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

	fmt.Println(req)
	query, args, err := u.db.Sq.Builder.Insert(u.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.tableName, "create"))
	}

	fmt.Println("3")


	fmt.Println(query)
	fmt.Println(args...)

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		fmt.Println(err, "<--------")
		return nil, u.db.Error(err)
	}


	return req, nil
}

func (u *userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
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
	data := map[string]any{
		"first_name": req.FirstName, 
		"last_name": req.LastName, 
		"username": req.Username, 
		"bio": req.Bio,
		"website": req.Website,
		"is_active": req.IsActive,
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