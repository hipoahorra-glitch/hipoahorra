document.addEventListener('DOMContentLoaded', () => {
    const formatNumber = (value) => new Intl.NumberFormat('es-ES', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(value);
    const calculatorForm = document.querySelector('form[action="/calcular#resultados"]');

    if (!calculatorForm) {
        return;
    }

    const inputs = calculatorForm.querySelectorAll('input, select');
    const resultNetSavings = document.querySelector('[data-live-net-savings]');
    const resultVerdict = document.querySelector('[data-live-verdict]');
    const resultCard = document.querySelector('[data-live-result-card]');
    const contactLink = document.querySelector('[data-live-contact]');
    const switchingSavingsElement = document.querySelector('[data-live-switching-savings]');
    const switchingMonthly = document.querySelector('[data-live-switching-monthly]');
    const detailBasePayment = document.querySelector('[data-live-base-payment]');
    const detailDiscountedPayment = document.querySelector('[data-live-discounted-payment]');
    const detailDiscountedRate = document.querySelector('[data-live-discounted-rate]');
    const detailTotalInterestBase = document.querySelector('[data-live-total-interest-base]');
    const detailTotalInterestDiscounted = document.querySelector('[data-live-total-interest-discounted]');
    const detailInterestSavings = document.querySelector('[data-live-interest-savings]');
    const detailMonthlySavings = document.querySelector('[data-live-monthly-savings]');
    const detailAnnualSavings = document.querySelector('[data-live-annual-savings]');
    const detailBankMonthly = document.querySelector('[data-live-bank-monthly]');
    const detailBankAnnual = document.querySelector('[data-live-bank-annual]');
    const detailExternalMonthly = document.querySelector('[data-live-external-monthly]');
    const detailExternalAnnual = document.querySelector('[data-live-external-annual]');
    const detailInsuranceDifference = document.querySelector('[data-live-insurance-difference]');

    const updateResults = () => {
        const formData = new FormData(calculatorForm);

        fetch('/api/calcular', {
            method: 'POST',
            headers: { 'Accept': 'application/json' },
            body: new URLSearchParams(formData)
        })
            .then((response) => response.json())
            .then((payload) => {
                if (!payload?.result) {
                    return;
                }

                const switchingSavings = payload.result.annualSwitchingSavings ?? 0;
                const isPositive = switchingSavings > 0;
                if (resultVerdict) {
                    resultVerdict.textContent = isPositive ? 'SÍ COMPENSA' : 'NO COMPENSA';
                    resultVerdict.classList.toggle('text-emerald-700', isPositive);
                    resultVerdict.classList.toggle('text-red-600', !isPositive);
                }
                if (resultNetSavings) {
                    resultNetSavings.classList.toggle('text-emerald-700', isPositive);
                    resultNetSavings.classList.toggle('text-red-600', !isPositive);
                }
                if (resultCard) {
                    resultCard.classList.toggle('border-emerald-200', isPositive);
                    resultCard.classList.toggle('border-red-200', !isPositive);
                    resultCard.classList.toggle('from-emerald-50', isPositive);
                    resultCard.classList.toggle('to-emerald-100/60', isPositive);
                    resultCard.classList.toggle('from-red-50', !isPositive);
                    resultCard.classList.toggle('to-rose-100/60', !isPositive);
                }
                if (contactLink) {
                    contactLink.classList.toggle('hidden', !isPositive);
                }
                if (switchingSavingsElement) {
                    switchingSavingsElement.textContent = `${switchingSavings > 0 ? '+' : ''}${formatNumber(switchingSavings)}`;
                }
                if (switchingMonthly) {
                    switchingMonthly.textContent = `${switchingSavings > 0 ? '+' : ''}${formatNumber(switchingSavings / 12)}`;
                }
                if (detailBasePayment) {
                    detailBasePayment.textContent = `${formatNumber(payload.result.monthlyPaymentWithoutBonuses ?? 0)} €/mes`;
                }
                if (detailDiscountedPayment) {
                    detailDiscountedPayment.textContent = `${formatNumber(payload.result.monthlyPaymentWithBonuses ?? 0)} €/mes`;
                }
                if (detailDiscountedRate) {
                    detailDiscountedRate.textContent = `${formatNumber(payload.result.discountedRate ?? 0)}%`;
                }
                if (detailTotalInterestBase) {
                    detailTotalInterestBase.textContent = `${formatNumber(payload.result.totalInterestWithoutBonuses ?? 0)} €`;
                }
                if (detailTotalInterestDiscounted) {
                    detailTotalInterestDiscounted.textContent = `${formatNumber(payload.result.totalInterestWithBonuses ?? 0)} €`;
                }
                if (detailInterestSavings) {
                    detailInterestSavings.textContent = formatNumber(payload.result.mortgageInterestSavings ?? 0);
                }
                if (detailMonthlySavings) {
                    detailMonthlySavings.textContent = `${formatNumber(payload.result.monthlySavings ?? 0)} €/mes`;
                }
                if (detailAnnualSavings) {
                    detailAnnualSavings.textContent = `${formatNumber(payload.result.annualSavings ?? 0)} €/año`;
                }
                if (detailBankMonthly) {
                    detailBankMonthly.textContent = `${formatNumber(payload.result.bankInsuranceMonthly ?? 0)} €/mes`;
                }
                if (detailBankAnnual) {
                    detailBankAnnual.textContent = `${formatNumber(payload.result.bankInsuranceAnnual ?? 0)} €/año`;
                }
                if (detailExternalMonthly) {
                    detailExternalMonthly.textContent = `${formatNumber(payload.result.externalInsuranceMonthly ?? 0)} €/mes`;
                }
                if (detailExternalAnnual) {
                    detailExternalAnnual.textContent = `${formatNumber(payload.result.externalInsuranceAnnual ?? 0)} €/año`;
                }
                if (detailInsuranceDifference) {
                    detailInsuranceDifference.textContent = formatNumber(Math.abs(payload.result.insuranceAnnualDifference ?? 0));
                }
            })
            .catch(() => { });
    };

    inputs.forEach((input) => {
        input.addEventListener('input', updateResults);
        input.addEventListener('change', updateResults);
    });

    updateResults();
});
