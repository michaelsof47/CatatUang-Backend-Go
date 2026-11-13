FROM golang:1.26-alpine AS builder
WORKDIR /src

#Cache dependency layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/catatuang

# Runtime Stage
FROM alpine:3.21

#ca-certificates dibutuhkan untuk firebase-admin-sdk 
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/catatuang /app/catatuang

#Folder Upload mount sebagai volume
RUN mkdir -p /app/uploads && chown -R app:app /app
USER app

EXPOSE 8080
CMD ["/app/catatuang"]
