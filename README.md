# Golang Search Engine

A prototype search engine with a web crawler and full text search through UI and API integrations written in Go.

## Core features (MVP)

- Web crawler
- Custom full text search engine
- Cronjobs
- Admin UI to control web crawler
- JSON API with JWT auth and rate limiting

## Tech stack

- **Framework:** Fiber
- **Database:** PostgreSQL
- **ORM:** Gorm
- **Auth:** custom JWT tokens through Cookie
- **UI:** Templ, TailwindCSS
- **Deployment:** Docker, Railway