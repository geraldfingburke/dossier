# Build stage for Go backend
FROM golang:1.24 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY server/ ./server/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./server/cmd/main.go

# Build stage for Vue frontend
FROM node:20 AS frontend-builder
WORKDIR /app
COPY client/package*.json ./
RUN npm ci
COPY client/ ./
RUN npm run build

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/bin/server .

# Copy frontend build (optional, can be served separately)
COPY --from=frontend-builder /app/dist ./client/dist

EXPOSE 8080

CMD ["./server"]
