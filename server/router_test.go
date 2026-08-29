package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHomeRouteRendersSuccessfully(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	NewRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Comparador de ahorro hipotecario") {
		t.Fatalf("expected rendered page title in response body")
	}
}

func TestDataProtectionPolicyRouteRendersSuccessfully(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/politica-proteccion-datos", nil)
	w := httptest.NewRecorder()

	NewRouter().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Política de protección de datos") {
		t.Fatalf("expected privacy policy title in response body")
	}
}
