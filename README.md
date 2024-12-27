# Golang Search Engine <!-- omit in toc -->

A prototype search engine with a web crawler and full text search through UI and API integrations written in Go.

## Table of contents <!-- omit in toc -->

- [Core features (MVP)](#core-features-mvp)
- [Tech stack](#tech-stack)
- [Getting started](#getting-started)
  - [Running locally](#running-locally)
  - [Login credentials](#login-credentials)

## Core features (MVP)

- [Web crawler](https://en.wikipedia.org/wiki/Web_crawler)
- Custom [full text search engine](https://www.mongodb.com/resources/basics/full-text-search)
- Cronjobs
- Admin UI to configure the web crawler
- JSON API with JWT auth and [rate limiting](https://www.cloudflare.com/learning/bots/what-is-rate-limiting/#:~:text=Rate%20limiting%20is%20a%20strategy,kinds%20of%20malicious%20bot%20activity.)

## Tech stack

- **Go framework:** [Fiber](https://gofiber.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **ORM:** [Gorm](https://gorm.io/)
- **Auth:** JWT tokens using Cookies
- **UI:** [Templ](https://templ.guide/), [HTMX](https://htmx.org/), [DaisyUI](https://daisyui.com/), [TailwindCSS](https://tailwindcss.com/)
- **Deployment:** [Docker](https://www.docker.com/), [Air](https://github.com/air-verse/air), [Railway](https://railway.com/)

## Getting started

### Running locally

1. Clone the repo: `git clone https://github.com/sebastian-nunez/golang-search-engine`
2. Install [`Go`](https://go.dev/doc/install) 1.23 or greater
3. Install [`Docker Desktop`](https://docs.docker.com/compose/install/)
4. Install [`Air`](https://github.com/air-verse/air) : `go install github.com/air-verse/air@latest`
5. Run docker compose: `docker compose up`
6. Run the app: `air`
7. Open the crawler settings dashboard: [http://localhost:3000/](http://localhost:3000/) (see the `Login credentials` section below for instructions on how to login)

### Login credentials

For testing purposes, some dummy admin credentials can be inserted into the database.

Simply hit the following API endpoint: `GET /api/v1/create-admin` (or just open in the browser [http://localhost:3000/api/v1/create-admin](http://localhost:3000/api/v1/create-admin))

- **Email:** `jdoe@google.com`
- **Password:** `password`

<!-- TODO: add API docs for the /search endpoints -->
