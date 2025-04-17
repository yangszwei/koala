# Koala Frontend

This directory contains the frontend code for the Koala application. It's designed to be embedded into the Go backend using Go's `embed` package. In production, the frontend is served directly from the Go binary, eliminating the need for a separate frontend server.

## 🛠️ Development Setup

### Prerequisites

- [Node.js](https://nodejs.org/) version 20
- [npm](https://www.npmjs.com/)

### Installation

1. Navigate to the `/web` directory:

   ```bash
   cd web
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

### Configuration

You can configure the backend API target by creating a `.env` file in the `/web` directory.

Example:

```
VITE_BACKEND_TARGET=http://localhost:8080
```

To get started, copy the example file:

```bash
cp .env.example .env
```

### Running the Development Server

Start the development server with hot reloading:

```bash
npm run dev
```

The application will be available at http://localhost:5173/ by default.

## 📦 Building for Production

To build the frontend assets for production:

```bash
npm run build
```

This will generate the production-ready assets in the `dist/` directory.  
These assets are intended to be embedded into the Go backend.

## 🔗 Integration with Go Backend

The `dist/` directory is embedded into the Go application using the `embed` package.  
Ensure that you rebuild the frontend whenever you make changes, so the latest assets are included in the Go binary.

## 🧑‍💻 Code Style

To ensure consistency and maintainability, all contributors should follow the established code style:

- Use the provided Prettier and ESLint setup for formatting and linting.
- Run `npm run format` to automatically fix formatting issues.
- Run `npm run lint` to catch potential code issues.~~

It's recommended to configure your editor (e.g., VSCode, WebStorm) to format on save and use the project’s eslint/prettier settings.
