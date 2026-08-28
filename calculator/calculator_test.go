package calculator

import "testing"

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
