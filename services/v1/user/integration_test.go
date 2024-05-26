package userv1

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	dbtest "github.com/farolinar/dealls-bumble/internal/common/db/test"
	"github.com/farolinar/dealls-bumble/internal/common/jwt"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// integration testing for creating user
func TestUser_Integration_CreateUserSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()
	url := "/v1/user/register"
	var resp UserAuthenticationResponse

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
	userService := NewService(cfg, userRepo)
	userHandler := NewHandler(cfg, userService)

	// serviceData, err := userService.Create(ctx, getUserCreatePayload())
	// assert.NoError(t, err)
	// assert.NotEmpty(t, serviceData.Token)

	user := getUserCreatePayload()
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Error encoding JSON: %v", err)
	}
	payload := string(jsonData)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	userHandler.CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("error unmarshaling resp json: %v", err)
	}
	assert.Equal(t, servicebase.CodeSuccess, resp.Code)
	assert.NotEmpty(t, resp.Data.Token)

	// verify token is valid
	_, err = jwt.VerifyAndGetSubject(cfg.App.Secret, resp.Data.Token)
	assert.NoError(t, err)
}

func TestUser_Integration_CreateUserErrAlreadyExists(t *testing.T) {
	ctx := context.Background()
	cfg := getConfig()

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
	userService := NewService(cfg, userRepo)

	_, err = userService.Create(ctx, getUserCreatePayload())
	assert.NoError(t, err)

	_, err = userService.Create(ctx, getUserCreatePayload())
	assert.ErrorIs(t, err, ErrAlreadyExists)
}
