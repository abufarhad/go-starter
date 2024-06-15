package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/abufarhad/golang-starter-rest-api/domain"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/user/controller"
	"github.com/stretchr/testify/assert"
)

type MockUserService struct {
}

func (m *MockUserService) CreateUser(user dto.UserReq) (*domain.User, *errors.RestErr) {
	if user.Email == "existing@example.com" {
		return nil, errors.NewAlreadyExistError(errors.ErrInvalidEmailOrPassword)
	}

	// Mocking the hashed password generation
	hashedPassword := "mocked_hashed_password"

	return &domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: &hashedPassword,
	}, nil
}

func (m *MockUserService) UpdateUser(userID int, user dto.UpdateUserReq) (*domain.User, *errors.RestErr) {
	if userID == 123 {
		// Simulate updating user details
		updatedUser := &domain.User{
			Name: user.Name,
		}
		return updatedUser, nil
	}
	return nil, errors.NewNotFoundError(errors.ErrInvalidEmailOrPassword)
}

func TestCreateUser(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// Setup router and controller
	router := gin.Default()
	mockUserService := &MockUserService{}
	controller.NewUserController(router.Group(""), mockUserService)

	testCases := []struct {
		name           string
		requestPayload string
		expectedStatus int
	}{
		{
			name:           "Valid Input",
			requestPayload: `{"email": "valid@example.com", "password": "pass123"}`,
			expectedStatus: http.StatusOK,
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
		{
			name:           "Unique Email",
			requestPayload: `{"email": "new@example.com", "password": "pass123"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Duplicate Email",
			requestPayload: `{"email": "existing@example.com", "password": "password123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid Input - Password Hashed",
			requestPayload: `{"email": "valid@example.com", "password": "pass123"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			requestBody := []byte(testcase.requestPayload)

			req, err := http.NewRequest("POST", "/v1/user/signup", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, testcase.expectedStatus, w.Code)

		})
	}

}
