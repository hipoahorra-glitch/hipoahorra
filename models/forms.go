package models

type CalculationForm struct {
	MortgageAmount     float64
	Years              int
	AnnualInterest     float64
	Age                int
	HomeInsuranceBonus float64
	LifeInsuranceBonus float64
	OtherBonuses       float64
	Coverage           string
}

type ContactForm struct {
	Name          string
	Email         string
	Phone         string
	PendingAmount string
	Message       string
}
