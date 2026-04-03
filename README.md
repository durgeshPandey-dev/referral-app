# Referral Email Automation System

## Features

- Upload Excel with HR emails
- Bulk email sending
- Worker pool for concurrency
- Structured logging
- SendGrid integration

## Workspace

This project uses [Nx](https://nx.dev) — a build system with monorepo support,
caching, and task orchestration.

### Apps

| App                 | Description                                                                            |
| ------------------- | -------------------------------------------------------------------------------------- |
| `api`               | Go backend — handles file uploads, queues bulk email jobs, and dispatches via SendGrid |
| `@org/frontend`     | React/Vite frontend — UI for uploading HR email lists and tracking referral campaigns  |
| `@org/frontend-e2e` | Playwright end-to-end tests for the frontend                                           |

### Running the apps

```bash
make run PROJECT=api TARGET=serve           # Start the API server
make run PROJECT=@org/frontend TARGET=serve # Start the frontend dev server
```

Run `make help` to see all available projects, targets, and commands.

## Setup

1. Clone repo
2. Create configs/.env
3. Install dependencies:

```bash
  npm install
  go mod tidy
```

## API

POST /upload  
GET /health

## Roadmap

referral might transition into
a full-fledged job management portal for candidates with:

- Resume review using LLM,
- Keyword extraction from resume,
- Job scraping off internet based on keywords,
- Cold reachout for jobs with thread based tracking,
...and much more.

Development roadmap can be inferenced from [issue tracker](docs/TODO.md).

## Docs

There are no API or dashboard docs available (read required) right now.
APIs have been listed [above](#api).
In case dedicated docs are added in future,
docs will be available inrespective directories
for [api](docs/api/), [dashboard](docs/dashboard/) and [design](docs/design/).
