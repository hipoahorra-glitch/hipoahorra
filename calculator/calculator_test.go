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
