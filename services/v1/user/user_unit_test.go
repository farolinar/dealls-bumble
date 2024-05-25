package userv1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/farolinar/dealls-bumble/internal/common/password"
	servicebase "github.com/farolinar/dealls-bumble/services/base"
	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestUser_Unit_CreateUser(t *testing.T) {
	var validationErrorBody UserAuthenticationResponse
	validationErrorBody.Code = servicebase.Code4XX

	type fields struct {
		svc func() Service
	}
	type args struct {
		rw http.ResponseWriter
		r  func() *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		respBody   UserAuthenticationResponse
		httpStatus int
	}{
		{
			name: "Validation name error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Name = "a"
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Validation email error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Email = "notemail"
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Validation username error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Username = "a"
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Validation password error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Password = "short"
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Validation sex error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Sex = ""
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Validation birthdate error - returns 400",
			fields: fields{
				svc: func() Service {
					db, _, _ := sqlmock.New()
					userRepo := NewRepository(db)
					mockUserService := NewService(userRepo)

					return mockUserService
				},
			},
			args: args{
				r: func() *http.Request {
					user := getUserCreatePayload()
					user.Birthdate = "23-10-1999"
					jsonData, err := json.Marshal(user)
					if err != nil {
						t.Fatalf("Error encoding JSON: %v", err)
					}
					payload := string(jsonData)
					req, err := http.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(payload))
					if err != nil {
						t.Fatal(err)
					}

					return req
				},
			},
			respBody:   validationErrorBody,
			httpStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Handler{
				service: tt.fields.svc(),
			}

			requestRecorder := httptest.NewRecorder()
			c.CreateUser(requestRecorder, tt.args.r())
			var resp UserAuthenticationResponse
			err := json.NewDecoder(requestRecorder.Body).Decode(&resp)
			if err != nil {
				t.Fatalf("Error decoding JSON: %v", err)
				return
			}
			assert.Equal(t, tt.httpStatus, requestRecorder.Code)
			assert.Equal(t, tt.respBody.Code, resp.Code)
		})
	}
}

func getUserCreatePayload() UserCreatePayload {
	return UserCreatePayload{
		Name:       "Tav",
		Email:      "example@email.com",
		Username:   "tavishere",
		Password:   "Pass12345!",
		Sex:        "female",
		Birthdate:  "1999-10-23",
		TimeLayout: "2006-01-02",
	}
}

func getTestUserEntity() (user User, err error) {
	payload := getUserCreatePayload()

	hashedPassword, err := password.Hash(payload.Password)
	if err != nil {
		return
	}

	dateString := "2023-10-23"
	layout := "2006-01-02"

	birthdate, err := time.Parse(layout, payload.Birthdate)
	if err != nil {
		return
	}
	currentDate, err := time.Parse(layout, dateString)
	if err != nil {
		return
	}

	return User{
		UID:            "uid123",
		Name:           payload.Name,
		Email:          payload.Email,
		Username:       payload.Username,
		HashedPassword: &hashedPassword,
		Sex:            payload.Sex,
		Birthdate:      birthdate,
		CreatedAt:      currentDate,
	}, nil
}
