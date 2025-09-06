FROM golang:1.25 AS builder
LABEL authors="Griffin Skudder"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o modbustohttp .

FROM gcr.io/distroless/static-debian12
LABEL authors="Griffin Skudder"
WORKDIR /app
COPY --from=builder /app/modbustohttp .

ENTRYPOINT ["./modbustohttp"]