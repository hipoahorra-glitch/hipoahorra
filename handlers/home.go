package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mortgage-bonus-calculator/calculator"
	"mortgage-bonus-calculator/models"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	form := defaultCalculationForm()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "Comparador de ahorro hipotecario",
		"form":        form,
		"showResults": false,
	})
}

func Calculate(c *gin.Context) {
	form := parseCalculationForm(c)
	result := buildResult(form)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "Comparador de ahorro hipotecario",
		"result":      result,
		"form":        form,
		"showResults": true,
	})
}

func CalculateJSON(c *gin.Context) {
	form := parseCalculationForm(c)
	result := buildResult(form)

	c.JSON(http.StatusOK, gin.H{
		"result": result,
		"form":   form,
	})
}

func Contact(c *gin.Context) {
	form := models.ContactForm{
		Name:          strings.TrimSpace(c.PostForm("name")),
		Email:         strings.TrimSpace(c.PostForm("email")),
		Phone:         strings.TrimSpace(c.PostForm("phone")),
		PendingAmount: strings.TrimSpace(c.PostForm("pendingAmount")),
		Message:       strings.TrimSpace(c.PostForm("message")),
	}

	log.Printf("Nueva solicitud de contacto: %+v", form)
	fmt.Println("Formulario de contacto enviado")
	c.Redirect(http.StatusSeeOther, "/")
}

func defaultCalculationForm() models.CalculationForm {
	return models.CalculationForm{
		MortgageAmount:       300000,
		Years:                20,
		AnnualInterest:       3.0,
		Age:                  55,
		BankInsuranceMonthly: 350,
		HomeInsuranceBonus:   0.0,
		LifeInsuranceBonus:   0.15,
		OtherBonuses:         0.0,
		Coverage:             "death",
	}
}

func parseCalculationForm(c *gin.Context) models.CalculationForm {
	principal, _ := strconv.ParseFloat(c.PostForm("mortgageAmount"), 64)
	years, _ := strconv.Atoi(c.PostForm("years"))
	annualRate, _ := strconv.ParseFloat(c.PostForm("annualInterest"), 64)
	age, _ := strconv.Atoi(c.PostForm("age"))
	lifeBonus, _ := strconv.ParseFloat(c.PostForm("lifeInsuranceBonus"), 64)
	bankInsuranceMonthly, _ := strconv.ParseFloat(c.PostForm("bankInsuranceMonthly"), 64)

	return models.CalculationForm{
		MortgageAmount:       clampFloat(principal, 1000, 2000000),
		Years:                clampInt(years, 1, 40),
		AnnualInterest:       clampFloat(annualRate, 0.01, 15),
		Age:                  clampInt(age, 18, 85),
		BankInsuranceMonthly: clampFloat(bankInsuranceMonthly, 50, 1000),
		LifeInsuranceBonus:   clampFloat(lifeBonus, 0, 5),
	}
}

func clampFloat(value, minimum, maximum float64) float64 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func clampInt(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func buildResult(form models.CalculationForm) calculator.Result {
	return calculator.CalculateScenario(
		form.MortgageAmount,
		form.Years,
		form.AnnualInterest,
		form.Age,
		form.LifeInsuranceBonus,
		form.BankInsuranceMonthly,
		calculator.ExternalInsurancePremium(form.MortgageAmount, form.Age),
	)
}
