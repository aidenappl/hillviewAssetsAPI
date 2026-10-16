# hillviewAssetsAPI

Equipment asset tracking and checkouts for the Hillview TV studio.

> **Hillview platform** · Go API · `/assets/v1.1` on api.hillview.tv

---

## Overview

Equipment tracking for the Hillview TV studio: assets, the students/staff who borrow them, and checkout/check-in. Used by the team dashboard. It has **no authentication of its own** — treat every route as public.

## Role in the Hillview ecosystem

One of five Go APIs behind `api.hillview.tv` (core, video, auth, links, assets), consumed by hillview.tv and team.hillview.tv. Reports to Monitor zone **hillview** as `hillview-assets-api`.

## Tech stack

Go 1.26 · gorilla/mux · squirrel · MariaDB · go-monitor · Docker · Lattice

## Getting started

### Prerequisites

Go 1.26+, the `dev` CLI, access to a MariaDB with the Hillview schema.

### Setup

Export the variables listed in [AGENTS.md §4](AGENTS.md#4-running-building--testing), then `dev`.

## Development

| Command | What it does |
|---|---|
| `dev` | Run the API |
| `dev build` | Build `bin/app` |
| `dev test` | Run tests |
| `dev fmt` | gofmt |
| `dev vet` | go vet |
| `dev check` | fmt + vet + test |
| `dev tidy` | go mod tidy |

## Project structure

See [AGENTS.md §3](AGENTS.md#3-project-structure).

## Deployment

Push to `main` → GitHub Actions builds `registry.appleby.cloud/hillview-assets-api:latest` → Lattice stack 19 redeploys the container.

## Contributing & further reading

- [AGENTS.md](AGENTS.md) — conventions, telemetry, operations
- [docs/API.md](docs/API.md) — endpoint reference
- Related: hillviewCoreAPI, hillviewVideoAPI, hillviewAuthAPI, hillviewLinksAPI, hillviewAssetsAPI, hillview-tv-web, hillview-team-web, hillview-live-web
