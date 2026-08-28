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
	result := buildResult(form)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "Comparador de ahorro hipotecario",
		"result":      result,
		"form":        form,
		"showResults": true,
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
		MortgageAmount:     150000,
		Years:              25,
		AnnualInterest:     3.0,
		Age:                40,
		HomeInsuranceBonus: 0.0,
		LifeInsuranceBonus: 0.3,
		OtherBonuses:       0.0,
		Coverage:           "death",
	}
}

func parseCalculationForm(c *gin.Context) models.CalculationForm {
	principal, _ := strconv.ParseFloat(c.PostForm("mortgageAmount"), 64)
	years, _ := strconv.Atoi(c.PostForm("years"))
	annualRate, _ := strconv.ParseFloat(c.PostForm("annualInterest"), 64)
	age, _ := strconv.Atoi(c.PostForm("age"))
	homeBonus, _ := strconv.ParseFloat(c.PostForm("homeInsuranceBonus"), 64)
	lifeBonus, _ := strconv.ParseFloat(c.PostForm("lifeInsuranceBonus"), 64)
	otherBonus, _ := strconv.ParseFloat(c.PostForm("otherBonuses"), 64)
	coverage := strings.TrimSpace(c.PostForm("coverage"))
	if coverage == "" {
		coverage = "death"
	}

	return models.CalculationForm{
		MortgageAmount:     principal,
		Years:              years,
		AnnualInterest:     annualRate,
		Age:                age,
		HomeInsuranceBonus: homeBonus,
		LifeInsuranceBonus: lifeBonus,
		OtherBonuses:       otherBonus,
		Coverage:           coverage,
	}
}

func buildResult(form models.CalculationForm) calculator.Result {
	return calculator.CalculateCase(
		form.MortgageAmount,
		form.Years,
		form.AnnualInterest,
		form.Age,
		form.HomeInsuranceBonus,
		form.LifeInsuranceBonus,
		form.OtherBonuses,
		estimateBankPremium(form.MortgageAmount, form.Age, form.Coverage),
		estimateExternalPremium(form.MortgageAmount, form.Age, form.Coverage),
		form.Coverage,
	)
}

func estimateBankPremium(capital float64, age int, coverage string) float64 {
	premium := calculator.EstimatePremium(capital, age, coverage)
	if premium <= 0 {
		return 45
	}
	return premium
}

func estimateExternalPremium(capital float64, age int, coverage string) float64 {
	premium := estimateBankPremium(capital, age, coverage) * 0.6
	if premium < 15 {
		return 15
	}
	return premium
}
