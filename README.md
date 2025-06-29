# Cars Viewer

This repository provides a small Cars API (under `api/`) and a Go based web viewer.

## Running the server

1. Install Go 1.20 or later.
2. Build or run the server:

```bash
go run ./
```

The server listens on `http://localhost:8080` and serves the web UI.

Images for car models are served from `/images/`.

## API Endpoints

The Go server exposes the following endpoints:

- `GET /api/models` – list car models with optional query parameters `search`, `manufacturerId`, `categoryId`.
- `GET /api/models/{id}` – details for a specific model.
- `GET /api/models/compare?ids=1,2` – compare multiple models.
- `GET /api/manufacturers` and `GET /api/manufacturers/{id}` – manufacturer info.
- `GET /api/categories` and `GET /api/categories/{id}` – category info.
- `GET /api/recommendations` – recommended models based on `manufacturerId` or `categoryId`.

## Development

Static files for the frontend live in the `web/` directory.
