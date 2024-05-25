package userv1

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, user *User) (err error)
	GetByUsername(ctx context.Context, username string) (user User, err error)
}

type dbRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, user *User) (err error) {
	q := `
        INSERT INTO users (uid, name, email, username, hashed_password, sex, birthdate)
        VALUES ($1, $2, $3, $4, $5, $6, $7);
    `
	_, err = d.db.ExecContext(ctx, q,
		user.UID, user.Name, user.Email, user.Username, user.HashedPassword, user.Sex,
		user.Birthdate)

	return
}

func (d dbRepository) GetByUsername(ctx context.Context, username string) (user User, err error) {
	q := `
        SELECT uid, name, email, username, hashed_password, sex, birthdate, created_at
        FROM users
        WHERE username = $1;
    `
	row := d.db.QueryRowContext(ctx, q, username)
	err = row.Scan(&user.UID, &user.Name, &user.Email, &user.Username, &user.HashedPassword,
		&user.Sex, &user.Birthdate, &user.CreatedAt)
	// if err == sql.ErrNoRows {
	//     return nil, ErrNotFound
	// }
	return
}
