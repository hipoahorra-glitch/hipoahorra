package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"unicode"

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
	log.Printf("Consulta de contacto recibida: Nombre=%q | Email=%q | Teléfono=%q | Importe pendiente=%q | Mensaje=%q", form.Name, form.Email, form.Phone, form.PendingAmount, form.Message)

	if err := sendContactEmail(form); err != nil {
		log.Printf("No se pudo enviar la consulta de contacto: %v", err)
		c.String(http.StatusInternalServerError, "No se pudo enviar la consulta. Inténtalo de nuevo más tarde.")
		return
	}

	log.Printf("Consulta de contacto enviada correctamente")
	c.Redirect(http.StatusSeeOther, "/")
}

func sendContactEmail(form models.ContactForm) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	username := strings.TrimSpace(os.Getenv("SMTP_USERNAME"))
	password := strings.Map(func(character rune) rune {
		if unicode.IsSpace(character) {
			return -1
		}
		return character
	}, os.Getenv("SMTP_PASSWORD"))
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if host == "" || port == "" || username == "" || password == "" || from == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}
	if strings.ContainsAny(host+port+username+from, "\r\n") {
		return fmt.Errorf("SMTP configuration contains an invalid line break")
	}
	log.Printf("SMTP preparado: Host=%q | Puerto=%q | Usuario=%q | Remitente=%q | Destinatario=%q", host, port, username, from, "hipoahorra@gmail.com")
	log.Printf("SMTP_PASSWORD recibido: longitud=%d caracteres", len(password))

	subjectName := strings.NewReplacer("\r", " ", "\n", " ").Replace(form.Name)
	subject := fmt.Sprintf("Nueva consulta de hipoteca – %s", subjectName)
	body := fmt.Sprintf("Me gustaría recibir más información sobre mis opciones hipotecarias.\n\nNombre: %s\n\nEmail: %s\n\nTeléfono: %s\n\nImporte pendiente de hipoteca: %s\n\nMensaje:\n\n%s\n", form.Name, form.Email, form.Phone, form.PendingAmount, form.Message)
	message := strings.Join([]string{
		"From: " + from,
		"To: hipoahorra@gmail.com",
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(host+":"+port, auth, from, []string{"hipoahorra@gmail.com"}, []byte(message))
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
