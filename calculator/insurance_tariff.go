package calculator

// ExternalInsurancePremium returns the monthly premium from the NN death-only tariff.
// Values are interpolated by capital and age, matching the reference matrix.
func ExternalInsurancePremium(capital float64, age int) float64 {
	return nnFallecimientoPremium(capital, age)
}

func interpolateCapital(capital float64, capitals, premiums []float64) float64 {
	if capital <= capitals[0] {
		return premiums[0]
	}
	for index := 1; index < len(capitals); index++ {
		if capital <= capitals[index] {
			fraction := (capital - capitals[index-1]) / (capitals[index] - capitals[index-1])
			return premiums[index-1] + (premiums[index]-premiums[index-1])*fraction
		}
	}
	return premiums[len(premiums)-1]
}
