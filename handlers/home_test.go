package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCalculateJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	form := url.Values{}
	form.Set("mortgageAmount", "150000")
	form.Set("years", "25")
	form.Set("annualInterest", "3.0")
	form.Set("age", "40")
	form.Set("homeInsuranceBonus", "0.0")
	form.Set("lifeInsuranceBonus", "0.3")
	form.Set("otherBonuses", "0.0")
	form.Set("coverage", "death")

	req := httptest.NewRequest(http.MethodPost, "/api/calcular", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	CalculateJSON(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload struct {
		Result struct {
			MonthlySavings float64 `json:"monthlySavings"`
			NetSavings     float64 `json:"netSavings"`
		} `json:"result"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload.Result.MonthlySavings <= 0 {
		t.Fatalf("expected positive monthly savings in response, got %.2f", payload.Result.MonthlySavings)
	}

	if payload.Result.NetSavings == 0 {
		t.Fatalf("expected non-zero net savings in response")
	}
}
