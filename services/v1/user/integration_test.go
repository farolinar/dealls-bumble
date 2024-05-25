package userv1

import (
	"context"
	"database/sql"
	"testing"

	dbtest "github.com/farolinar/dealls-bumble/internal/common/db/test"
	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// integration testing for creating user
func TestUser_Integration_CreateUserSuccess(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := dbtest.CreatePostgresContainer(ctx)
	if err != nil {
		t.Fatalf("error creating postgres container: %v", err)
	}

	db, err := sql.Open("pgx", pgContainer.ConnectionString)
	if err != nil {
		t.Fatalf("unable to connect to database: %v\n", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close db: %v", err)
		}
	})

	userRepo := NewRepository(db)
	userService := NewService(userRepo)

	resp, err := userService.Create(ctx, getUserCreatePayload())
	assert.NoError(t, err)

	assert.NotEmpty(t, resp.Token)
}

func TestUser_Integration_CreateUserErrAlreadyExists(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := dbtest.CreatePostgresContainer(ctx)
	if err != nil {
		t.Fatalf("error creating postgres container: %v", err)
	}

	db, err := sql.Open("pgx", pgContainer.ConnectionString)
	if err != nil {
		t.Fatalf("unable to connect to database: %v\n", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close db: %v", err)
		}
	})

	userRepo := NewRepository(db)
	userService := NewService(userRepo)

	_, err = userService.Create(ctx, getUserCreatePayload())
	assert.NoError(t, err)

	_, err = userService.Create(ctx, getUserCreatePayload())
	assert.ErrorIs(t, err, ErrAlreadyExists)
}
