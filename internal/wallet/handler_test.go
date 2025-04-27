package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"budget_manager/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

type mockUserHandler struct{}

func (m *mockUserHandler) Register(ctx *gin.Context) {}
func (m *mockUserHandler) Login(ctx *gin.Context)    {}
func (m *mockUserHandler) Logout(ctx *gin.Context)   {}

type mockSessionManager struct{}

func (m *mockSessionManager) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func TestCreateWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockService(ctrl)
	h := NewHandler(s)
	router := router.SetupRouter(h, &mockUserHandler{}, &mockSessionManager{})

	expectedWallet := Wallet{
		ID:         1,
		UserID:     1,
		Title:      "test",
		General:    10000,
		Operations: make([]Operation, 0),
	}

	// Success
	s.EXPECT().Save(expectedWallet).Return(expectedWallet, nil)
	successData, err := json.Marshal(&expectedWallet)
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodPost,
		"/wallet/create",
		bytes.NewReader(successData),
	)
	if err != nil {
		t.Fatalf("failed to create a request: %v", err)
	}

	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("wrong status: got %v, want %v", status, http.StatusOK)
	}

	b, err := json.Marshal(&expectedWallet)
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}
	if want := string(b); w.Body.String() != want {
		t.Errorf("wrong response: got %v, want %v", w.Body.String(), want)
	}

	// Error
	s.EXPECT().
		Save(expectedWallet).
		Return(Wallet{}, fmt.Errorf("no results"))
	brokenData := successData

	w = httptest.NewRecorder()

	req, err = http.NewRequest(
		http.MethodPost,
		"/wallet/create",
		bytes.NewReader(brokenData),
	)
	if err != nil {
		t.Fatalf("failed to create a request: %v", err)
	}

	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("wrong status: got %v, want %v", status, http.StatusBadRequest)
	}

	errMsg := "message: bad request"
	if w.Body.String() != errMsg {
		t.Errorf("wrong response: got %v, want %s", w.Body.String(), errMsg)
	}
}

func TestShowWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockService(ctrl)
	h := NewHandler(s)
	router := router.SetupRouter(h, &mockUserHandler{}, &mockSessionManager{})

	var userID int64 = 1
	wallet := Wallet{
		ID:         1,
		UserID:     userID,
		Title:      "test",
		General:    10000,
		Operations: make([]Operation, 0),
	}

	// Success
	s.EXPECT().ShowWallet(userID).Return(wallet, nil)

	b, err := json.Marshal(&map[string]any{"user_id": userID})
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodGet,
		"/wallet/show",
		bytes.NewReader(b),
	)
	if err != nil {
		t.Fatalf("failed to create a request: %v", err)
	}

	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status: got %v, want %v", status, http.StatusOK)
	}

	successData, err := json.Marshal(&wallet)
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}
	if want := string(successData); w.Body.String() != want {
		t.Errorf("wrong response: got %v, want %v", w.Body.String(), want)
	}
}

func TestAddOperation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := NewMockService(ctrl)
	h := NewHandler(s)
	router := router.SetupRouter(h, &mockUserHandler{}, &mockSessionManager{})

	op := Operation{
		ID:       1,
		WalletID: 1,
		Type:     "income",
		Amount:   1000,
	}

	// Success
	s.EXPECT().AddOperation(int64(1), op).Return(nil)
	successData, err := json.Marshal(&OperationOptions{
		UserID:    1,
		Operation: op,
	})
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}

	w := httptest.NewRecorder()

	req, err := http.NewRequest(
		http.MethodPost,
		"/wallet/operation/add",
		bytes.NewReader(successData),
	)
	if err != nil {
		t.Fatalf("failed to create a request: %v", err)
	}

	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Wrong status: got %v, want %v", status, http.StatusOK)
	}

	if want := "message: Success!"; w.Body.String() != want {
		t.Errorf("wrong response: got %v, want %v", w.Body.String(), want)
	}
}
