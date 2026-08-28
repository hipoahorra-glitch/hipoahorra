package calculator

import "math"

// Result contains the calculated values for the mortgage bonus comparison.
type Result struct {
	MortgageAmount               float64 `json:"mortgageAmount"`
	Years                        int     `json:"years"`
	AnnualInterest               float64 `json:"annualInterest"`
	Age                          int     `json:"age"`
	HomeInsuranceBonus           float64 `json:"homeInsuranceBonus"`
	LifeInsuranceBonus           float64 `json:"lifeInsuranceBonus"`
	OtherBonuses                 float64 `json:"otherBonuses"`
	Coverage                     string  `json:"coverage"`
	MonthlyPaymentWithoutBonuses float64 `json:"monthlyPaymentWithoutBonuses"`
	MonthlyPaymentWithBonuses    float64 `json:"monthlyPaymentWithBonuses"`
	MonthlySavings               float64 `json:"monthlySavings"`
	AnnualSavings                float64 `json:"annualSavings"`
	BankInsurancePremium         float64 `json:"bankInsurancePremium"`
	ExternalInsurancePremium     float64 `json:"externalInsurancePremium"`
	SavingsChangingInsurer       float64 `json:"savingsChangingInsurer"`
	NetSavings                   float64 `json:"netSavings"`
	DiscountedRate               float64 `json:"discountedRate"`
}

// NetSavingsForDisplay returns the savings amount as a positive magnitude for presentation.
func (r Result) NetSavingsForDisplay() float64 {
	if r.NetSavings < 0 {
		return -r.NetSavings
	}
	return r.NetSavings
}

// MortgagePayment calculates the monthly payment for a standard amortizing mortgage.
func MortgagePayment(principal float64, annualRatePct float64, years int) float64 {
	n := float64(years * 12)
	r := annualRatePct / 100 / 12
	if r == 0 {
		return principal / n
	}
	return principal * r / (1 - math.Pow(1+r, -n))
}

// CalculateCase computes the mortgage and insurance comparison.
func CalculateCase(principal float64, years int, annualRatePct float64, age int, bonHome float64, bonLife float64, bonOther float64, bankPremium float64, insurerPremium float64, coverage string) Result {
	basePayment := MortgagePayment(principal, annualRatePct, years)
	totalBonus := bonHome + bonLife + bonOther
	discountedRate := math.Max(annualRatePct-totalBonus, 0)
	discountedPayment := MortgagePayment(principal, discountedRate, years)

	monthlySavings := basePayment - discountedPayment
	annualSavings := monthlySavings * 12

	bankInsurancePremium := bankPremium
	externalInsurancePremium := insurerPremium
	premiumSaving := bankInsurancePremium - externalInsurancePremium
	netSavings := premiumSaving - monthlySavings

	if coverage == "death-plus-disability" {
		bankInsurancePremium *= 1.12
		externalInsurancePremium *= 1.1
		premiumSaving = bankInsurancePremium - externalInsurancePremium
		netSavings = premiumSaving - monthlySavings
	}

	return Result{
		MortgageAmount:               principal,
		Years:                        years,
		AnnualInterest:               annualRatePct,
		Age:                          age,
		HomeInsuranceBonus:           bonHome,
		LifeInsuranceBonus:           bonLife,
		OtherBonuses:                 bonOther,
		Coverage:                     coverage,
		MonthlyPaymentWithoutBonuses: basePayment,
		MonthlyPaymentWithBonuses:    discountedPayment,
		MonthlySavings:               monthlySavings,
		AnnualSavings:                annualSavings,
		BankInsurancePremium:         bankInsurancePremium,
		ExternalInsurancePremium:     externalInsurancePremium,
		SavingsChangingInsurer:       premiumSaving,
		NetSavings:                   netSavings,
		DiscountedRate:               discountedRate,
	}
}

// EstimatePremium creates an illustrative premium estimate based on capital and age.
func EstimatePremium(capital float64, age int, coverage string) float64 {
	base := 0.00085 * capital
	if age >= 60 {
		base *= 1.08
	} else if age >= 45 {
		base *= 1.03
	}

	if coverage == "death-plus-disability" {
		base *= 1.16
	}

	return math.Round(base/12*100) / 100
}
