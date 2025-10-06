# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
ENV PORT=8081
COPY --from=build /bin/app /app/app
USER nonroot:nonroot
EXPOSE 8081
ENTRYPOINT ["/app/app"]
