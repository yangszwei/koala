### Frontend builder ###
FROM node:20-alpine AS web-builder

WORKDIR /app/web

# Install deps separately for better cache hit
COPY web/package.json web/package-lock.json ./
RUN npm install

# Copy full source & build
COPY web/ ./
RUN npm run build


### Backend builder ###
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Copy go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy entire backend codebase
COPY . ./

# Copy built frontend into backend (e.g., for embedding or serving)
COPY --from=web-builder /app/web/dist ./web/dist

# Build the Go binary
RUN go build -o koala ./cmd/koala


### Final runtime image ###
FROM alpine:latest

WORKDIR /app

# Copy built binary
COPY --from=go-builder /app/koala .

EXPOSE 8000

ENTRYPOINT ["./koala"]
