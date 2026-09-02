document.addEventListener('DOMContentLoaded', () => {
    const formatNumber = (value) => new Intl.NumberFormat('es-ES', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(value);
    const calculatorForm = document.querySelector('form[action="/calcular#resultados"]');
    const contactForm = document.querySelector('[data-contact-form]');
    const contactStatus = contactForm?.querySelector('[data-contact-status]');

    contactForm?.addEventListener('submit', async (event) => {
        event.preventDefault();
        const submitButton = contactForm.querySelector('button[type="submit"]');
        submitButton.disabled = true;
        try {
            const response = await fetch(contactForm.action, {
                method: 'POST',
                headers: { Accept: 'application/json' },
                body: new URLSearchParams(new FormData(contactForm))
            });
            const payload = await response.json();
            if (!response.ok) throw new Error(payload.message || 'No se pudo enviar la consulta.');
            contactStatus.textContent = payload.message;
            contactStatus.className = 'rounded-xl bg-emerald-50 p-3 text-sm font-semibold text-emerald-700';
            window.dataLayer = window.dataLayer || [];
            window.dataLayer.push({
                event: 'contact_form_success',
                form_name: 'contacto'
            });
            if (typeof window.gtag_report_conversion === 'function') {
                window.gtag_report_conversion();
            }
            contactForm.reset();
        } catch (error) {
            contactStatus.textContent = error.message || 'No se pudo enviar la consulta. Inténtalo de nuevo más tarde.';
            contactStatus.className = 'rounded-xl bg-red-50 p-3 text-sm font-semibold text-red-700';
        } finally {
            submitButton.disabled = false;
        }
    });

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
    const switchingLabel = document.querySelector('[data-live-switching-label]');
    const switchingMessage = document.querySelector('[data-live-switching-message]');
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
    const insuranceChart = document.querySelector('[data-insurance-chart]');
    const insuranceChartCaption = document.querySelector('[data-insurance-chart-caption]');
    const insuranceChartGross = document.querySelector('[data-insurance-chart-gross]');
    const insuranceChartLostBonus = document.querySelector('[data-insurance-chart-lost-bonus]');
    const insuranceChartTotal = document.querySelector('[data-insurance-chart-total]');
    const insuranceChartYears = document.querySelector('[data-insurance-chart-years]');
    const insuranceChartInsight = document.querySelector('[data-insurance-chart-insight]');
    const insuranceChartTooltip = insuranceChart ? document.createElement('div') : null;
    const insuranceChartState = { interactivePoints: [], ratio: 1 };

    const formatWholeEuro = (value) => new Intl.NumberFormat('es-ES', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
    }).format(value);

    if (insuranceChart && insuranceChartTooltip) {
        insuranceChart.parentElement.style.position = 'relative';
        insuranceChartTooltip.className = 'pointer-events-none absolute hidden min-w-[220px] rounded-2xl border border-slate-200 bg-white/95 px-4 py-3 text-sm text-slate-700 shadow-[0_18px_50px_rgba(15,23,42,0.16)] backdrop-blur';
        insuranceChart.parentElement.appendChild(insuranceChartTooltip);
    }

    const drawInsuranceChart = () => {
        if (!insuranceChart) return;
        const points = JSON.parse(insuranceChart.dataset.chartValues || '[]');
        const context = insuranceChart.getContext('2d');
        const ratio = window.devicePixelRatio || 1;
        const width = insuranceChart.clientWidth * ratio;
        const height = insuranceChart.clientHeight * ratio;
        const padding = { top: 28 * ratio, right: 18 * ratio, bottom: 34 * ratio, left: 56 * ratio };
        const plotWidth = width - padding.left - padding.right;
        const plotHeight = height - padding.top - padding.bottom;
        insuranceChart.width = width;
        insuranceChart.height = height;
        context.clearRect(0, 0, width, height);
        insuranceChartState.interactivePoints = [];
        insuranceChartState.ratio = ratio;
        if (!points.length) return;
        const values = points.flatMap((point) => [point.bankAnnual, point.externalAnnual]);
        const minimum = Math.min(...values, 0);
        const maximum = Math.max(...values, 1);
        const range = Math.max(maximum - minimum, 1);
        const paddedMinimum = Math.max(0, minimum - range * 0.08);
        const paddedMaximum = maximum + range * 0.16;
        const scaledRange = Math.max(paddedMaximum - paddedMinimum, 1);
        const x = (index) => padding.left + index * (plotWidth / Math.max(points.length - 1, 1));
        const y = (value) => padding.top + (1 - (value - paddedMinimum) / scaledRange) * plotHeight;
        const drawGrid = () => {
            context.save();
            context.strokeStyle = 'rgba(148, 163, 184, 0.24)';
            context.lineWidth = 1 * ratio;
            context.fillStyle = '#64748b';
            context.font = `${11 * ratio}px sans-serif`;
            context.textAlign = 'right';
            context.textBaseline = 'middle';
            for (let step = 0; step <= 4; step += 1) {
                const value = paddedMinimum + (scaledRange * step / 4);
                const yPos = y(value);
                context.beginPath();
                context.moveTo(padding.left, yPos);
                context.lineTo(width - padding.right, yPos);
                context.stroke();
                context.fillText(`${formatWholeEuro(value)}€`, padding.left - (8 * ratio), yPos);
            }
            context.restore();
        };
        const drawGap = () => {
            context.save();
            for (let index = 0; index < points.length; index += 1) {
                const bankY = y(points[index].bankAnnual);
                const externalY = y(points[index].externalAnnual);
                context.strokeStyle = points[index].annualSaving >= 0 ? 'rgba(5, 150, 105, 0.24)' : 'rgba(220, 38, 38, 0.18)';
                context.lineWidth = 3 * ratio;
                context.beginPath();
                context.moveTo(x(index), bankY);
                context.lineTo(x(index), externalY);
                context.stroke();
            }
            context.restore();
        };
        const drawLine = (key, color, fillColor) => {
            context.beginPath();
            context.strokeStyle = color;
            context.lineWidth = 3 * ratio;
            context.lineJoin = 'round';
            context.lineCap = 'round';
            points.forEach((point, index) => index ? context.lineTo(x(index), y(point[key])) : context.moveTo(x(index), y(point[key])));
            context.stroke();
            context.fillStyle = fillColor;
            points.forEach((point, index) => {
                const xPos = x(index);
                const yPos = y(point[key]);
                context.beginPath();
                context.arc(xPos, yPos, 3.5 * ratio, 0, Math.PI * 2);
                context.fill();
                insuranceChartState.interactivePoints.push({
                    age: point.age,
                    bankAnnual: point.bankAnnual,
                    externalAnnual: point.externalAnnual,
                    annualSaving: point.annualSaving,
                    series: key === 'bankAnnual' ? 'Banco' : 'Externo',
                    x: xPos / ratio,
                    y: yPos / ratio
                });
            });
        };
        const drawXAxis = () => {
            context.save();
            context.strokeStyle = 'rgba(100, 116, 139, 0.55)';
            context.lineWidth = 1.5 * ratio;
            context.beginPath();
            context.moveTo(padding.left, height - padding.bottom);
            context.lineTo(width - padding.right, height - padding.bottom);
            context.stroke();
            context.fillStyle = '#475569';
            context.font = `${11 * ratio}px sans-serif`;
            context.textAlign = 'center';
            context.textBaseline = 'top';
            const tickIndexes = new Set([0, Math.floor((points.length - 1) / 2), points.length - 1]);
            tickIndexes.forEach((index) => {
                const xPos = x(index);
                context.beginPath();
                context.moveTo(xPos, height - padding.bottom);
                context.lineTo(xPos, height - padding.bottom + (6 * ratio));
                context.stroke();
                context.fillText(String(points[index].age), xPos, height - padding.bottom + (8 * ratio));
            });
            context.restore();
        };

        const drawCurrentAgeMarker = () => {
            const currentIndex = points.findIndex((point) => point.isCurrentAge);
            if (currentIndex === -1) return;
            const xPos = x(currentIndex);
            const bankY = y(points[currentIndex].bankAnnual);
            const externalY = y(points[currentIndex].externalAnnual);

            context.save();
            context.strokeStyle = 'rgba(15, 23, 42, 0.24)';
            context.setLineDash([6 * ratio, 6 * ratio]);
            context.lineWidth = 2 * ratio;
            context.beginPath();
            context.moveTo(xPos, padding.top);
            context.lineTo(xPos, height - padding.bottom);
            context.stroke();
            context.setLineDash([]);

            context.fillStyle = '#0f172a';
            context.font = `${11 * ratio}px sans-serif`;
            context.textAlign = 'center';
            context.textBaseline = 'bottom';
            context.fillText(`Edad ${points[currentIndex].age}`, xPos, Math.min(bankY, externalY) - (10 * ratio));

            context.fillStyle = '#ffffff';
            context.strokeStyle = '#0f172a';
            context.lineWidth = 2 * ratio;
            [bankY, externalY].forEach((yPos) => {
                context.beginPath();
                context.arc(xPos, yPos, 5 * ratio, 0, Math.PI * 2);
                context.fill();
                context.stroke();
            });
            context.restore();
        };

        drawGrid();
        drawGap();
        drawLine('bankAnnual', '#4c7b86', '#4c7b86');
        drawLine('externalAnnual', '#f59e0b', '#f59e0b');
        drawXAxis();
        drawCurrentAgeMarker();

        if (insuranceChartInsight) {
            const currentPoint = points.find((point) => point.isCurrentAge) || points[0];
            const strongestGapPoint = points.reduce((best, point) => Math.abs(point.annualSaving) > Math.abs(best.annualSaving) ? point : best, points[0]);
            const currentDiffLabel = `${formatWholeEuro(Math.abs(currentPoint.annualSaving))} €/año`;
            if (currentPoint.annualSaving > 0) {
                insuranceChartInsight.textContent = `A los ${currentPoint.age} años, el banco sale ${currentDiffLabel} más caro que el externo. La mayor diferencia del tramo aparece a los ${strongestGapPoint.age} años.`;
            } else if (currentPoint.annualSaving < 0) {
                insuranceChartInsight.textContent = `A los ${currentPoint.age} años, el banco sale ${currentDiffLabel} más barato, pero la diferencia máxima del tramo aparece a los ${strongestGapPoint.age} años.`;
            } else {
                insuranceChartInsight.textContent = `A los ${currentPoint.age} años, ambos importes anuales están empatados. La diferencia más clara del tramo aparece a los ${strongestGapPoint.age} años.`;
            }
        }
    };

    const hideInsuranceTooltip = () => {
        if (!insuranceChartTooltip) return;
        insuranceChartTooltip.classList.add('hidden');
    };

    const showInsuranceTooltip = (event) => {
        if (!insuranceChart || !insuranceChartTooltip || !insuranceChartState.interactivePoints.length) return;
        const rect = insuranceChart.getBoundingClientRect();
        const pointerX = event.clientX - rect.left;
        const pointerY = event.clientY - rect.top;
        const threshold = 12;
        let nearestPoint = null;
        let nearestDistance = Infinity;

        insuranceChartState.interactivePoints.forEach((point) => {
            const distance = Math.hypot(pointerX - point.x, pointerY - point.y);
            if (distance < nearestDistance) {
                nearestDistance = distance;
                nearestPoint = point;
            }
        });

        if (!nearestPoint || nearestDistance > threshold) {
            hideInsuranceTooltip();
            return;
        }

        const savingLabel = `${nearestPoint.annualSaving >= 0 ? '+' : '-'}${formatNumber(Math.abs(nearestPoint.annualSaving))} €/año`;
        insuranceChartTooltip.innerHTML = `<p class="text-xs font-semibold uppercase tracking-[0.15em] text-slate-500">Edad ${nearestPoint.age}</p><p class="mt-2 text-sm font-semibold text-slate-900">${nearestPoint.series}</p><div class="mt-2 space-y-1"><p>Banco: <strong>${formatNumber(nearestPoint.bankAnnual)} €/año</strong></p><p>Externo: <strong>${formatNumber(nearestPoint.externalAnnual)} €/año</strong></p><p>Diferencia: <strong>${savingLabel}</strong></p></div>`;
        insuranceChartTooltip.classList.remove('hidden');

        const tooltipRect = insuranceChartTooltip.getBoundingClientRect();
        let left = nearestPoint.x + 14;
        let top = nearestPoint.y - tooltipRect.height - 14;
        if (left + tooltipRect.width > rect.width - 8) {
            left = nearestPoint.x - tooltipRect.width - 14;
        }
        if (top < 8) {
            top = nearestPoint.y + 14;
        }
        insuranceChartTooltip.style.left = `${left}px`;
        insuranceChartTooltip.style.top = `${top}px`;
    };

    drawInsuranceChart();
    window.addEventListener('resize', drawInsuranceChart);
    insuranceChart?.addEventListener('mousemove', showInsuranceTooltip);
    insuranceChart?.addEventListener('mouseleave', hideInsuranceTooltip);

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

                if (insuranceChart && payload.result.insuranceChartJSON) {
                    insuranceChart.dataset.chartValues = payload.result.insuranceChartJSON;
                    drawInsuranceChart();
                }
                if (insuranceChartCaption) {
                    insuranceChartCaption.textContent = `La gráfica es una estimación orientativa de primas anuales desde los ${payload.result.insuranceChartStartAge} hasta los ${payload.result.insuranceChartEndAge} años. Se resalta tu edad actual para ver el punto de partida real y, debajo, se separa el ahorro bruto del seguro de la bonificación hipotecaria que podrías perder.`;
                }
                if (insuranceChartGross) {
                    insuranceChartGross.textContent = `${formatNumber(payload.result.insuranceChartGrossSave ?? 0)} €`;
                }
                if (insuranceChartLostBonus) {
                    insuranceChartLostBonus.textContent = `${formatNumber(payload.result.insuranceChartLostBonus ?? 0)} €`;
                }
                if (insuranceChartTotal) {
                    insuranceChartTotal.textContent = `${formatNumber(payload.result.insuranceChartNetSave ?? 0)} €`;
                    insuranceChartTotal.classList.toggle('text-emerald-700', (payload.result.insuranceChartNetSave ?? 0) >= 0);
                    insuranceChartTotal.classList.toggle('text-rose-600', (payload.result.insuranceChartNetSave ?? 0) < 0);
                }
                if (insuranceChartYears) {
                    insuranceChartYears.textContent = String(payload.result.insuranceChartYears ?? 0);
                }

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
                if (switchingLabel) {
                    switchingLabel.textContent = isPositive ? 'Puedes ahorrar' : 'A dia de hoy perderias';
                }
                if (switchingSavingsElement) {
                    switchingSavingsElement.textContent = `${switchingSavings > 0 ? '+' : ''}${formatNumber(switchingSavings)}`;
                    switchingSavingsElement.classList.toggle('text-emerald-700', isPositive);
                    switchingSavingsElement.classList.toggle('text-red-600', !isPositive);
                }
                if (switchingMessage) {
                    switchingMessage.textContent = isPositive
                        ? 'Ahorro anual neto despues de tener en cuenta la perdida de la bonificacion.'
                        : 'Aun no merece la pena cambiarlo: la perdida de la bonificacion supera el ahorro del seguro externo.';
                }
                if (switchingMonthly) {
                    switchingMonthly.textContent = `${switchingSavings > 0 ? '+' : switchingSavings < 0 ? '-' : ''}${formatNumber(Math.abs(switchingSavings / 12))}`;
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
                    detailBankMonthly.textContent = `${formatNumber(payload.result.bankTariffMonthly ?? 0)} €/mes`;
                }
                if (detailBankAnnual) {
                    detailBankAnnual.textContent = `${formatNumber(payload.result.bankTariffAnnual ?? 0)} €/año`;
                }
                if (detailExternalMonthly) {
                    detailExternalMonthly.textContent = `${formatNumber(payload.result.externalInsuranceMonthly ?? 0)} €/mes`;
                }
                if (detailExternalAnnual) {
                    detailExternalAnnual.textContent = `${formatNumber(payload.result.externalInsuranceAnnual ?? 0)} €/año`;
                }
                if (detailInsuranceDifference) {
                    detailInsuranceDifference.textContent = formatNumber(Math.abs(payload.result.insuranceTariffDifference ?? 0));
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
