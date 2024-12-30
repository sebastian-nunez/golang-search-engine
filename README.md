# Golang Search Engine <!-- omit in toc -->

A prototype search engine with a web crawler and full text search through UI and API integrations written in Go.

## Table of contents <!-- omit in toc -->

- [Core features (MVP)](#core-features-mvp)
- [Tech stack](#tech-stack)
- [Getting started](#getting-started)
  - [Running locally](#running-locally)
  - [Login credentials](#login-credentials)
  - [Seeding URL(s) for the crawler](#seeding-urls-for-the-crawler)
- [API reference](#api-reference)
  - [`POST /api/v1/search`](#post-apiv1search)

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

While running the app locally, the main caveats are around seeding data (e.g. you will have to manually insert admin credentials to login as well and you will have to insert the first URL(s) for the crawler to start exploring).

### Running locally

1. Clone the repo: `git clone https://github.com/sebastian-nunez/golang-search-engine`
2. Install [`Go`](https://go.dev/doc/install) 1.23 or greater
3. Install [`Docker Desktop`](https://docs.docker.com/compose/install/)
4. Install [`Air`](https://github.com/air-verse/air) : `go install github.com/air-verse/air@latest`
5. Run docker compose: `docker compose up`
6. Run the app: `air`
7. You must open and log into the crawler settings dashboard: [http://localhost:3000/](http://localhost:3000/) (see the `Login credentials` section below for instructions on how to login). Otherwise, the initial crawler settings will NOT be set in the database.
8. You must manually seed the initial URL(s) for the crawler to begin exploring into the database (see the `Seeding URL(s) for the crawler` section below).

### Login credentials

For testing purposes, some dummy admin credentials can be inserted into the database.

Simply hit the following API endpoint: `GET /api/v1/create-admin` (or just open in the browser [http://localhost:3000/api/v1/create-admin](http://localhost:3000/api/v1/create-admin))

- **Email:** `jdoe@google.com`
- **Password:** `password`

<!-- TODO: add API docs for the /search endpoints -->

### Seeding URL(s) for the crawler

With a fresh database, the crawler will not have any websites to visit. So, we must give it a starting point(s).

For now, the only way is through an insert query into PostgreSQL:

```sql
-- id is a UUID (v4)
INSERT INTO crawled_pages (id, url)
VALUES ('532c357c-0397-4d9f-b63a-c06255ae747e', 'https://news.yahoo.com/');
```

TODO: add the API endpoint for this.

## API reference

To try out the API, I recommend downloading and using [Postman](https://www.postman.com/downloads/).

- **Base URL (prod):** TODO
- **Base URL (localhost):** `http://localhost:3000`

**Note:** To get the full URL, simply concatenate the Base URL to the API endpoint targeted (e.g. `http://localhost:3000/api/v1/search`)

### `POST /api/v1/search`

Allows users to search for indexed web pages which contain the query terms.

**JSON request body:**

```json
{
  "query": "sebastian nunez"
}
```

**JSON response:**

```json
{
  "results": [
    {
      "id": "b1a34930-b2a3-4f1f-9454-10e926754d55",
      "url": "https://www.sebastian-nunez.com/",
      "indexed": true,
      "success": true,
      "crawlDuration": 89205,
      "statusCode": 200,
      "lastTested": "2024-12-30T01:30:40.975935-05:00",
      "title": "Sebastian Nunez",
      "description": "Sebastian Nunez - Software Engineer Intern at Google | ex. UKG, JPMorgan Chase & Co. | Google Tech Exchange Scholar | Apple Pathways Alliance | Computer Science student at Florida International University",
      "headings": "Hey, I'm Sebastian!",
      "createdAt": "2024-12-30T01:31:34.671422-05:00",
      "updatedAt": "2024-12-30T01:31:34.671422-05:00",
      "deletedAt": null
    }
    // More results...
  ],
  "total": 1
}
```
