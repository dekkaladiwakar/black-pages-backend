FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build BOTH binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy BOTH binaries
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy migrations folder
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p uploads/resumes uploads/portfolios uploads/logos
EXPOSE 8080
CMD ["./main"]