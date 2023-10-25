FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /connectivity-wizard

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=build-stage /connectivity-wizard /app/connectivity-wizard
COPY static/ /app/static
COPY templates/ /app/templates

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/connectivity-wizard"]