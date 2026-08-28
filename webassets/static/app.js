document.addEventListener('DOMContentLoaded', () => {
    const calculatorForm = document.querySelector('form[action="/calcular#resultados"]');

    if (!calculatorForm) {
        return;
    }

    const inputs = calculatorForm.querySelectorAll('input, select');
    const resultPanel = document.querySelector('[data-live-result]');
    const resultTitle = resultPanel?.querySelector('[data-live-title]');
    const resultAmount = resultPanel?.querySelector('[data-live-amount]');
    const resultBadge = resultPanel?.querySelector('[data-live-badge]');
    const detailHomeBonus = document.querySelector('[data-live-bonus-home]');
    const detailLifeBonus = document.querySelector('[data-live-bonus-life]');
    const detailOtherBonus = document.querySelector('[data-live-bonus-other]');
    const detailCoverage = document.querySelector('[data-live-coverage]');
    const detailBasePayment = document.querySelector('[data-live-base-payment]');
    const detailDiscountedPayment = document.querySelector('[data-live-discounted-payment]');
    const detailMonthlySavings = document.querySelector('[data-live-monthly-savings]');
    const detailAnnualSavings = document.querySelector('[data-live-annual-savings]');

    const updateResults = () => {
        const formData = new FormData(calculatorForm);

        fetch('/api/calcular', {
            method: 'POST',
            headers: { 'Accept': 'application/json' },
            body: new URLSearchParams(formData)
        })
            .then((response) => response.json())
            .then((payload) => {
                if (!payload?.result || !resultPanel) {
                    return;
                }

                const savings = payload.result.netSavings ?? 0;
                const absoluteSavings = Math.abs(savings);
                const isPositive = savings > 0;

                if (resultTitle) {
                    resultTitle.textContent = isPositive ? 'Ahorro estimado' : 'Ahorro estimado';
                }

                if (resultAmount) {
                    resultAmount.textContent = `${absoluteSavings.toFixed(2)} €/mes`;
                }

                if (resultBadge) {
                    resultBadge.textContent = isPositive ? 'A favor del cambio' : 'Revisar con detalle';
                }

                if (detailHomeBonus) {
                    detailHomeBonus.textContent = `${(payload.result.homeInsuranceBonus ?? 0).toFixed(2)}%`;
                }
                if (detailLifeBonus) {
                    detailLifeBonus.textContent = `${(payload.result.lifeInsuranceBonus ?? 0).toFixed(2)}%`;
                }
                if (detailOtherBonus) {
                    detailOtherBonus.textContent = `${(payload.result.otherBonuses ?? 0).toFixed(2)}%`;
                }
                if (detailCoverage) {
                    const coverage = payload.result.coverage || payload.form.coverage || 'death';
                    detailCoverage.textContent = coverage === 'death-plus-disability' ? 'Muerte + invalidez permanente' : 'Muerte';
                }
                if (detailBasePayment) {
                    detailBasePayment.textContent = `${(payload.result.monthlyPaymentWithoutBonuses ?? 0).toFixed(2)} €/mes`;
                }
                if (detailDiscountedPayment) {
                    detailDiscountedPayment.textContent = `${(payload.result.monthlyPaymentWithBonuses ?? 0).toFixed(2)} €/mes`;
                }
                if (detailMonthlySavings) {
                    detailMonthlySavings.textContent = `${(payload.result.monthlySavings ?? 0).toFixed(2)} €/mes`;
                }
                if (detailAnnualSavings) {
                    detailAnnualSavings.textContent = `${(payload.result.annualSavings ?? 0).toFixed(2)} €/año`;
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
