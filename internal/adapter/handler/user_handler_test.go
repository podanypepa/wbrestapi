package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaveUC struct{ mock.Mock }

func (m *MockSaveUC) Execute(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockGetUC struct{ mock.Mock }

func (m *MockGetUC) Execute(externalID string) (*domain.User, error) {
	args := m.Called(externalID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func setupApp(saveUC port.SaveUserExecutor, getUC port.GetUserExecutor) *fiber.App {
	app := fiber.New()
	h := &UserHandler{SaveUC: saveUC, GetUC: getUC}
	h.RegisterRoutes(app)
	return app
}

func TestSaveUser_Success(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	user := domain.User{
		ExternalID:  "uuid-abc-123",
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	mockUC.On("Execute", mock.AnythingOfType("*domain.User")).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(respBody), "\"external_id\":\"uuid-abc-123\"")
}

func TestSaveUser_BadPayload(t *testing.T) {
	mockUC := new(MockSaveUC)
	app := setupApp(mockUC, nil)

	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetUser_Success(t *testing.T) {
	mockUC := new(MockGetUC)
	app := setupApp(nil, mockUC)

	externalID := "uuid-xyz-789"
	mockUser := &domain.User{
		ID:          1,
		ExternalID:  externalID,
		Name:        "Josef",
		Email:       "josef@example.com",
		DateOfBirth: time.Now(),
	}

	mockUC.On("Execute", externalID).Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/"+externalID, nil)
	resp, _ := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "\"external_id\":\""+externalID+"\"")
}

func TestGetUser_NotFound(t *testing.T) {
	mockUC := new(MockGetUC)
	app := setupApp(nil, mockUC)

	externalID := "not-found-uuid"
	mockUC.On("Execute", externalID).Return(&domain.User{}, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/"+externalID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
