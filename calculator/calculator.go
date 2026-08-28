package calculator

import (
	"encoding/json"
	"math"
)

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
	TotalInterestWithoutBonuses  float64 `json:"totalInterestWithoutBonuses"`
	TotalInterestWithBonuses     float64 `json:"totalInterestWithBonuses"`
	MortgageInterestSavings      float64 `json:"mortgageInterestSavings"`
	BankInsuranceMonthly         float64 `json:"bankInsuranceMonthly"`
	BankInsuranceAnnual          float64 `json:"bankInsuranceAnnual"`
	BankInsuranceTotal           float64 `json:"bankInsuranceTotal"`
	ExternalInsuranceMonthly     float64 `json:"externalInsuranceMonthly"`
	ExternalInsuranceAnnual      float64 `json:"externalInsuranceAnnual"`
	InsuranceAnnualDifference    float64 `json:"insuranceAnnualDifference"`
	AnnualLostBonusCost          float64 `json:"annualLostBonusCost"`
	AnnualSwitchingSavings       float64 `json:"annualSwitchingSavings"`
	SwitchingCompensates         bool    `json:"switchingCompensates"`
	NetSavingsOverLoan           float64 `json:"netSavingsOverLoan"`
	Compensates                  bool    `json:"compensates"`
	BankInsurancePremium         float64 `json:"bankInsurancePremium"`
	ExternalInsurancePremium     float64 `json:"externalInsurancePremium"`
	SavingsChangingInsurer       float64 `json:"savingsChangingInsurer"`
	NetSavings                   float64 `json:"netSavings"`
	DiscountedRate               float64 `json:"discountedRate"`
	BankTariffMonthly            float64 `json:"bankTariffMonthly"`
	InsuranceChartJSON           string  `json:"insuranceChartJSON"`
}

// NetSavingsForDisplay returns the savings amount as a positive magnitude for presentation.
func (r Result) NetSavingsForDisplay() float64 {
	if r.NetSavings < 0 {
		return -r.NetSavings
	}
	return r.NetSavings
}

func (r Result) InsuranceAnnualDifferenceForDisplay() float64 {
	return math.Abs(r.InsuranceAnnualDifference)
}

func (r Result) InsuranceMonthlyDifferenceForDisplay() float64 {
	return r.InsuranceAnnualDifferenceForDisplay() / 12
}

func (r Result) AnnualSwitchingSavingsForDisplay() float64 {
	return math.Abs(r.AnnualSwitchingSavings)
}

func (r Result) MonthlySwitchingSavingsForDisplay() float64 {
	return r.AnnualSwitchingSavingsForDisplay() / 12
}

func (r Result) NetSavingsOverLoanForDisplay() float64 {
	return math.Abs(r.NetSavingsOverLoan)
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

// MortgageInterest returns the total interest paid over the remaining loan term.
func MortgageInterest(principal float64, annualRatePct float64, years int) float64 {
	return MortgagePayment(principal, annualRatePct, years)*float64(years*12) - principal
}

type amortizationTotals struct {
	monthlyPayment float64
	totalInterest  float64
}

// amortizeTotals mirrors the reference app's month-by-month French schedule.
func amortizeTotals(principal, monthlyRate float64, months int) amortizationTotals {
	if principal <= 0 || months <= 0 {
		return amortizationTotals{}
	}

	monthlyPayment := principal / float64(months)
	if monthlyRate != 0 {
		monthlyPayment = principal * monthlyRate / (1 - math.Pow(1+monthlyRate, -float64(months)))
	}

	balance := principal
	totalInterest := 0.0
	for month := 1; month <= months; month++ {
		interest := balance * monthlyRate
		principalPayment := monthlyPayment - interest
		if month == months {
			principalPayment = balance
		}
		totalInterest += interest
		balance = math.Max(balance-principalPayment, 0)
	}

	return amortizationTotals{
		monthlyPayment: monthlyPayment,
		totalInterest:  totalInterest,
	}
}

// CalculateScenario calculates the mortgage, insurance, and final decision values.
func CalculateScenario(principal float64, years int, annualRatePct float64, age int, lifeBonus float64, bankMonthly float64, externalMonthly float64) Result {
	months := years * 12
	base := amortizeTotals(principal, annualRatePct/100/12, months)
	discountedRate := math.Max(annualRatePct-lifeBonus, 0)
	discounted := amortizeTotals(principal, discountedRate/100/12, months)
	basePayment := base.monthlyPayment
	discountedPayment := discounted.monthlyPayment
	baseInterest := base.totalInterest
	discountedInterest := discounted.totalInterest
	interestSavings := baseInterest - discountedInterest
	bankAnnual := bankMonthly * 12
	externalAnnual := externalMonthly * 12
	loanCost := bankAnnual * float64(years)
	netSavings := interestSavings - loanCost
	annualLostBonusCost := 0.0
	if years > 0 {
		annualLostBonusCost = interestSavings / float64(years)
	}
	annualSwitchingSavings := bankAnnual - externalAnnual - annualLostBonusCost

	return Result{
		MortgageAmount:               principal,
		Years:                        years,
		AnnualInterest:               annualRatePct,
		Age:                          age,
		LifeInsuranceBonus:           lifeBonus,
		MonthlyPaymentWithoutBonuses: basePayment,
		MonthlyPaymentWithBonuses:    discountedPayment,
		MonthlySavings:               basePayment - discountedPayment,
		AnnualSavings:                (basePayment - discountedPayment) * 12,
		TotalInterestWithoutBonuses:  baseInterest,
		TotalInterestWithBonuses:     discountedInterest,
		MortgageInterestSavings:      interestSavings,
		BankInsuranceMonthly:         bankMonthly,
		BankInsuranceAnnual:          bankAnnual,
		BankInsuranceTotal:           loanCost,
		ExternalInsuranceMonthly:     externalMonthly,
		ExternalInsuranceAnnual:      externalAnnual,
		InsuranceAnnualDifference:    bankAnnual - externalAnnual,
		AnnualLostBonusCost:          annualLostBonusCost,
		AnnualSwitchingSavings:       annualSwitchingSavings,
		SwitchingCompensates:         annualSwitchingSavings > 0,
		NetSavingsOverLoan:           netSavings,
		NetSavings:                   netSavings,
		DiscountedRate:               discountedRate,
		Compensates:                  netSavings > 0,
		BankTariffMonthly:            BankInsurancePremium(principal, age),
		InsuranceChartJSON:           insuranceChartJSON(principal),
	}
}

func insuranceChartJSON(capital float64) string {
	type point struct {
		Age      int     `json:"age"`
		Bank     float64 `json:"bank"`
		External float64 `json:"external"`
	}
	points := make([]point, 0, 30)
	for age := 30; age <= 59; age++ {
		points = append(points, point{Age: age, Bank: BankInsurancePremium(capital, age), External: nnFallecimientoPremium(capital, age)})
	}
	data, _ := json.Marshal(points)
	return string(data)
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
