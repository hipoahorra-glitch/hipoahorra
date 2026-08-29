# Mortgage Bonus Calculator

# hipoahorra

Small Go web app for comparing whether a bank's mortgage rate bonus offsets the cost of bundled insurance.

The interface is in Spanish and focuses on a simple homeowner workflow:

- enter mortgage details
- apply rate bonuses from home, life, and other products
- compare estimated monthly savings against estimated insurance costs
- review the result live in the form and in the detailed results section

## Overview

This project renders a server-side landing page with Gin templates and augments the calculator with a small client-side live update layer.

Main behaviors:



- `GET /` renders the full page with default calculation values
- `POST /calcular` renders the page again with updated results
- `POST /api/calcular` returns JSON for the live preview in the form
- `POST /contacto` logs contact form submissions and redirects back to `/`

## Stack

- Go `1.25`
- Gin
- HTML templates
- Tailwind via CDN
- Vanilla JavaScript

## Project Layout

```text
.
├── calculator/   core mortgage and premium calculations
├── handlers/     HTTP handlers for page, JSON, and contact form flows
├── models/       form models shared by handlers/templates
├── server/       Gin router setup
├── static/       browser JavaScript
└── templates/    page template and partials
```

Key files:

- `main.go`: app entrypoint, port selection, server startup
- `server/router.go`: static files, template registration, route wiring
- `handlers/home.go`: form parsing, HTML response, JSON response, contact handling
- `calculator/calculator.go`: payment, premium, and savings calculations
- `static/app.js`: live form updates through `/api/calcular`
- `templates/partials/`: page sections used by `templates/index.html`

## Requirements

- Go `1.25` or newer

Install dependencies:

```bash
go mod download
```

## Run Locally

Start the app:

```bash
go run .
```

The server listens on `http://localhost:8080` by default.

You can override the port:

```bash
PORT=3000 go run .
```

## Contact Email

The contact form sends submissions to `hipoahorra@gmail.com`. Configure SMTP before running the app:

```bash
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USERNAME=your-sender@gmail.com
export SMTP_PASSWORD=your-gmail-app-password
export SMTP_FROM=your-sender@gmail.com
go run .
```

For Gmail, use an App Password rather than the account password.

## Test

Run all tests:

```bash
go test ./...
```

Current automated coverage is focused on:

- core calculator behavior
- JSON calculation handler responses

## How The App Works

1. The initial page load uses default values from `handlers/home.go`.
2. The mortgage form posts to `/calcular` for a server-rendered result.
3. `static/app.js` also listens to form changes and sends the same inputs to `/api/calcular`.
4. The returned JSON updates the "Resumen en vivo" block without reloading the page.
5. Premium estimates are illustrative and are calculated in `calculator/calculator.go`.

## Development Notes

When editing the UI:

- update page composition in `templates/index.html`
- update form fields in `templates/partials/mortgage-form.html`
- update the detailed result presentation in `templates/partials/results-cards.html`
- keep `static/app.js` selectors in sync with any `data-live-*` attribute changes

When editing business logic:

- change mortgage or savings rules in `calculator/calculator.go`
- keep handler parsing in `handlers/home.go` aligned with any new form fields
- add or update tests before changing calculation semantics

## Notes And Limitations

- Insurance premiums are estimated, not insurer-backed quotes.
- Contact form submissions are logged to stdout; there is no persistence layer yet.
- Tailwind is loaded from a CDN in the page template, so no frontend build step is required.

## Next Improvements

- add validation feedback for incomplete or invalid form input
- persist contact submissions or send them to an email/CRM backend
- extend tests to cover full HTML form flows and edge-case calculations
