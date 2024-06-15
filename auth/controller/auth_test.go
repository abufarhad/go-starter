package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/abufarhad/golang-starter-rest-api/auth/controller"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/msgutil"
	"github.com/stretchr/testify/assert"
)

type MockAuthService struct {
}

func (m *MockAuthService) Login(cred *dto.LoginReq) (*dto.LoginResp, error) {
	if cred.Email == "valid@example.com" && cred.Password == "pass123" {
		return &dto.LoginResp{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}, nil
	}
	return nil, errors.ErrInvalidEmailOrPassword
}

func (m *MockAuthService) RefreshToken(refreshToken string) (*dto.LoginResp, error) {
	// Implement your logic here or return a default response for testing
	return nil, errors.ErrInvalidAccessToken
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	router := gin.Default()
	mockAuthService := &MockAuthService{}
	controller.NewAuthController(router.Group(""), mockAuthService)

	testCases := []struct {
		name           string
		requestPayload string
		expectedStatus int
	}{
		{
			name:           "Valid Credentials",
			requestPayload: `{"email": "valid@example.com", "password": "pass123"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Credentials - Wrong Password",
			requestPayload: `{"email": "valid@example.com", "password": "wrongpass"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Credentials - Wrong Email",
			requestPayload: `{"email": "invalid@example.com", "password": "pass123"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Input - Missing Email",
			requestPayload: `{"password": "pass123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Input - Missing Password",
			requestPayload: `{"email": "valid@example.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Input - Empty Request",
			requestPayload: `{}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Convert the request payload string to bytes
			requestBody := []byte(testCase.requestPayload)

			req, err := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatus, w.Code)

			// Check if the response body contains the password
			if strings.Contains(w.Body.String(), "pass123") {
				t.Errorf("response body contains password: %s", w.Body.String())
			}

			if testCase.expectedStatus == http.StatusOK {
				expectedResponse := msgutil.NewRestResp("auth response", &dto.LoginResp{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				})

				// Convert RestResp object to a JSON string
				expectedResponseJSON, err := json.Marshal(expectedResponse)
				if err != nil {
					t.Fatalf("error marshaling expected response: %v", err)
				}

				assert.JSONEq(t, string(expectedResponseJSON), w.Body.String())
			}
		})
	}
}
