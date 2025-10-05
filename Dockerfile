# Build the Vue SPA
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web .
RUN npm run build

# Build the Go backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Final runtime image
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app \
    && apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/server ./server
COPY --from=backend-builder /app/db ./db
COPY --from=frontend-builder /app/web/dist ./web/dist
ENV PORT=8080
EXPOSE 8080
USER app
ENTRYPOINT ["./server"]
