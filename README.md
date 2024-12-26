# Golang Search Engine <!-- omit in toc -->

A prototype search engine with a web crawler and full text search through UI and API integrations written in Go.

## Table of contents <!-- omit in toc -->

- [Core features (MVP)](#core-features-mvp)
- [Tech stack](#tech-stack)
- [Getting started](#getting-started)

## Core features (MVP)

- [Web crawler](https://en.wikipedia.org/wiki/Web_crawler)
- Custom [full text search engine](https://www.mongodb.com/resources/basics/full-text-search)
- Cronjobs
- Admin UI to monitor web crawler
- JSON API with JWT auth and [rate limiting](https://www.cloudflare.com/learning/bots/what-is-rate-limiting/#:~:text=Rate%20limiting%20is%20a%20strategy,kinds%20of%20malicious%20bot%20activity.)

## Tech stack

- **Framework:** [Fiber](https://gofiber.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **ORM:** [Gorm](https://gorm.io/)
- **Auth:** JWT tokens using Cookies
- **UI:** [Templ](https://templ.guide/), [HTMX](https://htmx.org/), [DaisyUI](https://daisyui.com/), [TailwindCSS](https://tailwindcss.com/)
- **Deployment:** [Docker](https://www.docker.com/), [Air](https://github.com/air-verse/air), [Railway](https://railway.com/)

## Getting started

1. Clone the repo: `git clone https://github.com/sebastian-nunez/golang-search-engine`
2. Install `Go` 1.23 or greater
3. Install `Air` : `go install github.com/air-verse/air@latest`
4. Run the app: `air`
