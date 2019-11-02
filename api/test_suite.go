package api

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/mocks"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type suite struct {
	gomock         *gomock.Controller
	userRepo       *mocks.MockUser
	accountService *mocks.MockAccount
	user           *models.User
	context        echo.Context
	request        *http.Request
	recorder       *httptest.ResponseRecorder
	api            *Api
}

func newTestSuite(t *testing.T, method string, body io.Reader, headers map[string]string) *suite {
	// Create mocks
	ctrl := gomock.NewController(t)

	mockUserRepository := mocks.NewMockUser(ctrl)
	mockAccountService := mocks.NewMockAccount(ctrl)

	// HTTP testing setup
	e := echo.New()
	req := httptest.NewRequest(method, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create logger instance
	logger := zap.NewNop().Sugar()

	// Create Api instance
	a := &Api{
		logger:         logger,
		accountService: mockAccountService,
		userRepo:       mockUserRepository,
	}

	return &suite{
		gomock:         ctrl,
		userRepo:       mockUserRepository,
		accountService: mockAccountService,
		request:        req,
		recorder:       rec,
		context:        c,
		api:            a,
	}
}

func (s *suite) authorize() {
	user := s.getDummyUser()

	s.user = user
	s.context.Set("user", user)
}

func (s *suite) getDummyUser() *models.User {
	ts, _ := time.Parse(time.RFC1123, time.RFC1123)

	// Authorizing user
	user := &models.User{
		ID: 1,
		Email: "test@test.com",
		Password: "my_hashed_password",
		Token: "",
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	return user
}

func (s *suite) close() {
	s.gomock.Finish()
}
