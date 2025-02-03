package auth_test

import (
	"adv-go/api/configs"
	"adv-go/api/internal/auth"
	"adv-go/api/internal/user"
	"adv-go/api/pkg/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	email = "a@1a.ru"
	password = "1234"
	name = "test user"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}

	return &handler, mock, nil
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password"}).AddRow(
		"a1@a.ru",
		"$2a$10$ruTbXRDNF9.ulyWTp1weUOl8W4WLyGYI6eehUe9gY6fyO8NYhebGe",
	)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	data, _ :=  json.Marshal(&auth.LoginRequest{
		Email: email,
		Password: password,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("status got: %d, want %d", w.Result().StatusCode, http.StatusOK)
	}
}


func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	data, _ :=  json.Marshal(&auth.RegisterRequest{
		Email: email,
		Password: password,
		Name: name,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("status got: %d, want %d", w.Result().StatusCode, http.StatusCreated)
	}
}