package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaveUC struct{ mock.Mock }

func (m *MockSaveUC) Execute(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockGetUC struct{ mock.Mock }

func (m *MockGetUC) Execute(ctx context.Context, externalID string) (*domain.User, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func setupApp(saveUC port.SaveUserExecutor, getUC port.GetUserExecutor) *fiber.App {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	v := validator.New()
	// Use JSON tag names for validation errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	h := &UserHandler{
		SaveUC:    saveUC,
		GetUC:     getUC,
		Logger:    logger,
		Validator: v,
	}
	h.RegisterRoutes(app)
	return app
}

func TestSaveUser_Success(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	reqPayload := UserRequest{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	mockUC.On("Execute", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "\"external_id\":\"550e8400-e29b-41d4-a716-446655440000\"")
}

func TestSaveUser_BadPayload(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSaveUser_ValidationError(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	// Invalid UUID should trigger validation error in handler
	reqPayload := UserRequest{
		ExternalID:  "invalid-uuid",
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "validation failed")
	assert.Contains(t, string(respBody), "external_id")
	assert.Contains(t, string(respBody), "must be a valid UUID")
}

func TestSaveUser_UserAlreadyExists(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	reqPayload := UserRequest{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	mockUC.On("Execute", mock.Anything, mock.AnythingOfType("*domain.User")).Return(domain.ErrUserAlreadyExists)

	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestGetUser_Success(t *testing.T) {
	mockUC := new(MockGetUC)
	app := setupApp(nil, mockUC)

	externalID := "550e8400-e29b-41d4-a716-446655440000"
	mockUser := &domain.User{
		ID:          1,
		ExternalID:  externalID,
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	mockUC.On("Execute", mock.Anything, externalID).Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/"+externalID, nil)
	resp, _ := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "\"external_id\":\""+externalID+"\"")
}

func TestGetUser_NotFound(t *testing.T) {
	mockUC := new(MockGetUC)
	app := setupApp(nil, mockUC)

	externalID := "550e8400-e29b-41d4-a716-446655440000"
	mockUC.On("Execute", mock.Anything, externalID).Return(nil, domain.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/"+externalID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGetUser_InternalError(t *testing.T) {
	mockUC := new(MockGetUC)
	app := setupApp(nil, mockUC)

	externalID := "550e8400-e29b-41d4-a716-446655440000"
	mockUC.On("Execute", mock.Anything, externalID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/"+externalID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
