FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOTOOLCHAIN=local
RUN CGO_ENABLED=0 GOOS=linux go build -o realworld-fiber-htmx .

FROM alpine:3

RUN apk add --no-cache wget

WORKDIR /app

COPY --from=builder /build/realworld-fiber-htmx .
COPY --from=builder /build/cmd/web/templates ./cmd/web/templates
COPY --from=builder /build/cmd/web/static ./cmd/web/static
COPY --from=builder /build/conduit.sqlite ./conduit.sqlite

ARG APP_PORT=8181
ENV APP_ADDR=0.0.0.0:${APP_PORT}
EXPOSE ${APP_PORT}

CMD ["./realworld-fiber-htmx"]