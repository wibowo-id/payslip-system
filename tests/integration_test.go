package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"payslip-system/internal/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin(t *testing.T) {
	r := app.SetupRouterForTest() // router dengan db in-memory

	// Register
	payload := map[string]string{
		"username": "testuser",
		"password": "secret123",
		"role":     "employee", // atau "admin", tergantung kebutuhan
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, 201, res.Code)

	// Login
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()

	r.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)

	var resp map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}
