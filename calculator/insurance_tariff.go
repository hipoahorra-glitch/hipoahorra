package calculator

// ExternalInsurancePremium returns the monthly premium from the NN death-only tariff.
// Values are interpolated by capital and age, matching the reference matrix.
func ExternalInsurancePremium(capital float64, age int) float64 {
	return nnFallecimientoPremium(capital, age)
}

// BankInsurancePremium returns the monthly premium from the ING reference tariff.
func BankInsurancePremium(capital float64, age int) float64 {
	return interpolateAgeCapital(capital, age, ingTariff)
}

var ingTariff = map[int][]float64{
	30: {9.01, 13.52, 18.34, 22.53, 27.09, 31.60, 35.83, 40.36, 45.00, 49.50, 53.33, 57.84, 62.34, 66.85, 70.82},
	40: {13.30, 19.95, 27.85, 33.24, 40.98, 47.96, 54.10, 60.74, 66.48, 73.13, 80.36, 86.89, 93.03, 99.66, 106.61},
	50: {18.20, 25.23, 32.26, 39.29, 46.32, 53.36, 60.39, 67.42, 74.45, 81.48, 88.51, 95.54, 102.58, 109.61, 116.64},
	55: {26.20, 36.83, 47.45, 58.00, 68.56, 79.11, 89.67, 100.25, 110.82, 121.40, 131.98, 142.56, 153.14, 163.71, 174.29},
	59: {31.45, 44.68, 57.91, 71.14, 84.37, 97.60, 110.83, 124.06, 137.30, 150.54, 163.77, 177.00, 190.23, 203.46, 216.69},
}

func interpolateAgeCapital(capital float64, age int, tariff map[int][]float64) float64 {
	ages := []int{30, 40, 50, 55, 59}
	if age <= ages[0] {
		return interpolateCapital(capital, tariffCapitals, tariff[ages[0]])
	}
	if age >= ages[len(ages)-1] {
		return interpolateCapital(capital, tariffCapitals, tariff[ages[len(ages)-1]])
	}
	for index := 1; index < len(ages); index++ {
		if age <= ages[index] {
			lower, upper := ages[index-1], ages[index]
			low := interpolateCapital(capital, tariffCapitals, tariff[lower])
			high := interpolateCapital(capital, tariffCapitals, tariff[upper])
			return low + (high-low)*float64(age-lower)/float64(upper-lower)
		}
	}
	return 0
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
