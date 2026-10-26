package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL    = "http://localhost:3001"
	serverPort = "3001"
)

var serverCmd *exec.Cmd

// TestMain runs integration tests with actual server
func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("PORT", serverPort)
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "users_test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSL", "disable")
	os.Setenv("RATE_LIMIT_MAX", "1000")

	// Start server
	serverCmd = exec.Command("go", "run", "main.go")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	if err := serverCmd.Start(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}

	// Wait for server to be ready
	if !waitForServer(baseURL+"/healthz", 30*time.Second) {
		serverCmd.Process.Kill()
		fmt.Println("Server failed to start in time")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if serverCmd != nil && serverCmd.Process != nil {
		serverCmd.Process.Kill()
	}

	os.Exit(code)
}

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func TestIntegration_HealthCheck(t *testing.T) {
	resp, err := http.Get(baseURL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
}

func TestIntegration_Metrics(t *testing.T) {
	resp, err := http.Get(baseURL + "/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Contains(t, string(body), "wbrestapi_http_requests_total")
}

func TestIntegration_SaveAndGetUser(t *testing.T) {
	// Create unique user
	userID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", time.Now().UnixNano()%1000000000000)
	payload := fmt.Sprintf(`{
		"external_id": "%s",
		"name": "Integration Test User",
		"email": "integration@test.com",
		"date_of_birth": "1990-01-01T00:00:00Z"
	}`, userID)

	// Save user
	resp, err := http.Post(baseURL+"/save", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var savedUser map[string]interface{}
	err = json.Unmarshal(body, &savedUser)
	require.NoError(t, err)

	assert.Equal(t, userID, savedUser["external_id"])
	assert.Equal(t, "Integration Test User", savedUser["name"])

	// Get user
	resp2, err := http.Get(baseURL + "/" + userID)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	body2, err := io.ReadAll(resp2.Body)
	require.NoError(t, err)

	var retrievedUser map[string]interface{}
	err = json.Unmarshal(body2, &retrievedUser)
	require.NoError(t, err)

	assert.Equal(t, userID, retrievedUser["external_id"])
	assert.Equal(t, "Integration Test User", retrievedUser["name"])
	assert.Equal(t, "integration@test.com", retrievedUser["email"])
}

func TestIntegration_SaveUser_ValidationError(t *testing.T) {
	payload := `{
		"external_id": "invalid-uuid",
		"name": "Test",
		"email": "test@test.com",
		"date_of_birth": "1990-01-01T00:00:00Z"
	}`

	resp, err := http.Post(baseURL+"/save", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_SaveUser_InvalidEmail(t *testing.T) {
	userID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", time.Now().UnixNano()%1000000000000)
	payload := fmt.Sprintf(`{
		"external_id": "%s",
		"name": "Test User",
		"email": "not-an-email",
		"date_of_birth": "1990-01-01T00:00:00Z"
	}`, userID)

	resp, err := http.Post(baseURL+"/save", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_GetUser_NotFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/550e8400-e29b-41d4-a716-999999999999")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_SaveUser_Duplicate(t *testing.T) {
	userID := fmt.Sprintf("550e8400-e29b-41d4-a716-%012d", time.Now().UnixNano()%1000000000000)
	payload := fmt.Sprintf(`{
		"external_id": "%s",
		"name": "Duplicate Test",
		"email": "duplicate@test.com",
		"date_of_birth": "1990-01-01T00:00:00Z"
	}`, userID)

	// First save
	resp1, err := http.Post(baseURL+"/save", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	// Second save with same external_id
	resp2, err := http.Post(baseURL+"/save", "application/json", strings.NewReader(payload))
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
}
