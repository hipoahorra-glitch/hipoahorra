package calculator

import (
	"encoding/json"
	"math"
	"testing"
)

func TestCalculateCase(t *testing.T) {
	result := CalculateCase(150000, 25, 3.0, 35, 0.0, 0.3, 0.0, 45.0, 25.0, "death")

	if result.MonthlyPaymentWithoutBonuses <= 0 {
		t.Fatalf("expected positive base payment, got %.2f", result.MonthlyPaymentWithoutBonuses)
	}

	if result.MonthlyPaymentWithBonuses >= result.MonthlyPaymentWithoutBonuses {
		t.Fatalf("expected discounted payment to be lower, got %.2f vs %.2f", result.MonthlyPaymentWithBonuses, result.MonthlyPaymentWithoutBonuses)
	}

	if result.MonthlySavings <= 0 {
		t.Fatalf("expected positive monthly savings, got %.2f", result.MonthlySavings)
	}

	if result.AnnualSavings != result.MonthlySavings*12 {
		t.Fatalf("expected annual savings to be monthly savings times 12")
	}

	if result.SavingsChangingInsurer <= 0 {
		t.Fatalf("expected positive premium saving, got %.2f", result.SavingsChangingInsurer)
	}
}

func TestNetSavingsForDisplay(t *testing.T) {
	positive := Result{NetSavings: 42.5}
	if got := positive.NetSavingsForDisplay(); got != 42.5 {
		t.Fatalf("expected positive savings magnitude, got %.2f", got)
	}

	negative := Result{NetSavings: -18.75}
	if got := negative.NetSavingsForDisplay(); got != 18.75 {
		t.Fatalf("expected positive magnitude for negative savings, got %.2f", got)
	}
}

func TestCalculateScenarioUsesBankInsuranceInput(t *testing.T) {
	result := CalculateScenario(150000, 20, 3.2, 40, 0.3, 35, 23)

	if result.BankInsuranceAnnual != 420 {
		t.Fatalf("expected bank insurance annual cost of 420, got %.2f", result.BankInsuranceAnnual)
	}

	expectedTotal := result.BankInsuranceAnnual * 20
	if result.BankInsuranceTotal != expectedTotal {
		t.Fatalf("expected bank insurance total of %.2f, got %.2f", expectedTotal, result.BankInsuranceTotal)
	}

	expectedNet := result.MortgageInterestSavings - expectedTotal
	if result.NetSavingsOverLoan != expectedNet {
		t.Fatalf("expected net savings of %.2f, got %.2f", expectedNet, result.NetSavingsOverLoan)
	}

	if result.Compensates != (expectedNet > 0) {
		t.Fatalf("compensation decision does not match net savings")
	}
}

func TestCalculateScenarioUsesFrenchAmortizationTotals(t *testing.T) {
	result := CalculateScenario(150000, 20, 3.2, 40, 0.3, 35, 23)

	expectedMonthly := MortgagePayment(150000, 3.2, 20)
	if result.MonthlyPaymentWithoutBonuses != expectedMonthly {
		t.Fatalf("expected French monthly payment %.10f, got %.10f", expectedMonthly, result.MonthlyPaymentWithoutBonuses)
	}

	if result.TotalInterestWithoutBonuses <= result.TotalInterestWithBonuses {
		t.Fatalf("expected bonified schedule to reduce total interest")
	}
	if result.MortgageInterestSavings != result.TotalInterestWithoutBonuses-result.TotalInterestWithBonuses {
		t.Fatalf("interest savings must equal the difference between both schedules")
	}
}

func TestExternalInsurancePremiumMatchesReferenceMatrix(t *testing.T) {
	if got := ExternalInsurancePremium(300000, 55); got != 131.98 {
		t.Fatalf("expected reference matrix premium of 131.98, got %.2f", got)
	}
}

func TestCalculateScenarioBuildsChartFromCurrentAge(t *testing.T) {
	result := CalculateScenario(300000, 20, 3.0, 55, 0.15, 350, ExternalInsurancePremium(300000, 55))

	if result.InsuranceChartStartAge != 30 {
		t.Fatalf("expected chart to start at age 30, got %d", result.InsuranceChartStartAge)
	}
	if result.InsuranceChartEndAge != 70 {
		t.Fatalf("expected chart to end at age 70, got %d", result.InsuranceChartEndAge)
	}
	if result.InsuranceChartYears != 41 {
		t.Fatalf("expected chart to span 41 years, got %d", result.InsuranceChartYears)
	}

	var points []struct {
		Age              int     `json:"age"`
		BankAnnual       float64 `json:"bankAnnual"`
		ExternalAnnual   float64 `json:"externalAnnual"`
		AnnualSaving     float64 `json:"annualSaving"`
		CumulativeSaving float64 `json:"cumulativeSaving"`
		IsCurrentAge     bool    `json:"isCurrentAge"`
	}
	if err := json.Unmarshal([]byte(result.InsuranceChartJSON), &points); err != nil {
		t.Fatalf("expected valid chart JSON, got error: %v", err)
	}
	if len(points) != 41 {
		t.Fatalf("expected 41 chart points, got %d", len(points))
	}
	if points[0].Age != 30 || points[len(points)-1].Age != 70 {
		t.Fatalf("expected chart ages 30..70, got %d..%d", points[0].Age, points[len(points)-1].Age)
	}
	if points[len(points)-1].BankAnnual <= points[len(points)-2].BankAnnual {
		t.Fatalf("expected projected bank annual premium to keep increasing after age 59")
	}
	if points[len(points)-1].ExternalAnnual <= points[len(points)-2].ExternalAnnual {
		t.Fatalf("expected projected external annual premium to keep increasing after age 59")
	}
	currentAgeCount := 0
	currentAgeIndex := -1
	for _, point := range points {
		if point.IsCurrentAge {
			currentAgeCount++
			currentAgeIndex = point.Age - points[0].Age
			if point.Age != 55 {
				t.Fatalf("expected current-age marker at 55, got %d", point.Age)
			}
		}
	}
	if currentAgeCount != 1 {
		t.Fatalf("expected exactly one current-age marker, got %d", currentAgeCount)
	}
	if currentAgeIndex <= 0 {
		t.Fatalf("expected current-age index to be greater than zero, got %d", currentAgeIndex)
	}
	expectedFromCurrentAge := points[len(points)-1].CumulativeSaving - points[currentAgeIndex-1].CumulativeSaving
	if math.Abs(result.InsuranceChartEstimatedSave-expectedFromCurrentAge) > 0.001 {
		t.Fatalf("expected estimated saving %.2f to match cumulative chart saving from current age %.2f", result.InsuranceChartEstimatedSave, expectedFromCurrentAge)
	}
}

func TestCalculateScenarioUsesTariffDifferenceForInsuranceChartSection(t *testing.T) {
	result := CalculateScenario(300000, 20, 3.0, 40, 0.15, 350, ExternalInsurancePremium(300000, 40))

	expectedAnnual := result.BankTariffMonthly * 12
	if result.BankTariffAnnual != expectedAnnual {
		t.Fatalf("expected bank tariff annual %.2f, got %.2f", expectedAnnual, result.BankTariffAnnual)
	}

	expectedDifference := result.BankTariffAnnual - result.ExternalInsuranceAnnual
	if result.InsuranceTariffDifference != expectedDifference {
		t.Fatalf("expected tariff difference %.2f, got %.2f", expectedDifference, result.InsuranceTariffDifference)
	}
}

func TestCalculateScenarioSplitsGrossLostAndNetChartEstimates(t *testing.T) {
	result := CalculateScenario(300000, 20, 3.0, 55, 0.15, 350, ExternalInsurancePremium(300000, 55))

	if result.InsuranceChartGrossSave != result.InsuranceChartEstimatedSave {
		t.Fatalf("expected gross chart save %.2f to match estimated save %.2f", result.InsuranceChartGrossSave, result.InsuranceChartEstimatedSave)
	}

	expectedLostBonus := result.AnnualLostBonusCost * float64(result.InsuranceChartEndAge-result.Age+1)
	if result.InsuranceChartLostBonus != expectedLostBonus {
		t.Fatalf("expected lost bonus %.2f, got %.2f", expectedLostBonus, result.InsuranceChartLostBonus)
	}

	expectedNet := result.InsuranceChartGrossSave - result.InsuranceChartLostBonus
	if result.InsuranceChartNetSave != expectedNet {
		t.Fatalf("expected net chart save %.2f, got %.2f", expectedNet, result.InsuranceChartNetSave)
	}
}

func TestCalculateScenarioLimitsChartToRemainingMortgageYears(t *testing.T) {
	result := CalculateScenario(300000, 3, 3.0, 40, 0.15, 350, ExternalInsurancePremium(300000, 40))

	if result.InsuranceChartStartAge != 30 || result.InsuranceChartEndAge != 42 {
		t.Fatalf("expected chart ages 30..42, got %d..%d", result.InsuranceChartStartAge, result.InsuranceChartEndAge)
	}
	if result.InsuranceChartYears != 13 {
		t.Fatalf("expected chart to span 13 years, got %d", result.InsuranceChartYears)
	}
}

func TestPremiumsProjectBeyondSourceMatrixToAge70(t *testing.T) {
	bankAt59 := BankInsurancePremium(300000, 59)
	bankAt70 := BankInsurancePremium(300000, 70)
	if bankAt70 <= bankAt59 {
		t.Fatalf("expected bank premium at 70 %.2f to exceed age 59 %.2f", bankAt70, bankAt59)
	}

	externalAt59 := ExternalInsurancePremium(300000, 59)
	externalAt70 := ExternalInsurancePremium(300000, 70)
	if externalAt70 <= externalAt59 {
		t.Fatalf("expected external premium at 70 %.2f to exceed age 59 %.2f", externalAt70, externalAt59)
	}
}
