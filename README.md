# Koala

A medical imaging and report search system.

---

Koala is a full-stack application for searching and organizing medical studies and reports. It leverages a Go backend and a modern Vite/React frontend, integrates FHIR and DICOMweb standards, and is powered by Elasticsearch for efficient querying.

## 🌟 Features

- 🔍 Search functionality for report content, modality, date range, patient details, and clinical categories
- 💡 Autocomplete support for faster and more accurate query entry
- 📄 Integration of FHIR DiagnosticReports with DICOMweb ImagingStudies
- ⚙️ Go backend providing structured APIs for search and data access
- 🖥️ Frontend built with Vite, React, and TailwindCSS for responsive user interfaces
- 🐳 Docker-ready build process for containerized deployment

## 📦 How to Build

### Build from Source

For users who prefer to build and host Koala manually without Docker:

#### Frontend

```bash
cd web
npm install
npm run build
```

This will generate static assets in the `web/dist` directory.

#### Backend

The backend requires the built frontend assets to compile successfully. Ensure you have the frontend built before proceeding.

```bash
go build -o koala cmd/koala/main.go
```

The compiled binary `koala` will be available in the project root.

### Docker Build

To build and run using Docker:

```bash
docker build -t koala .
```

This will create a Docker image named `koala`. You can run it with:

```bash
docker run -p 8080:8080 -v "./config.yaml:/app/config.yaml" koala
```

## ⚙️ Configuration

Create a config.yaml file in your project root for custom configurations. See the [`config/default.yaml`](config/default.yaml) for example configurations.

> **Note**: The `.env` file in `/web` is only used during frontend development to configure the Vite dev server. It does not affect backend configuration.

## 🛠️ Development Setup

### Prerequisites

- Go 1.21+
- Node.js 20+

### Clone the Repository

```bash
git clone https://github.com/yangszwei/koala.git
cd koala
```

### Frontend

See [`web/README.md`](web/README.md) for detailed setup, or run:

```bash
cd web
npm install
npm run dev
```

Koala frontend will be available at [http://localhost:5173](http://localhost:5173).

### Backend

```bash
go run cmd/koala/main.go
```

## 📜 License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
