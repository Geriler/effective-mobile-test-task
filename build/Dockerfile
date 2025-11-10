FROM golang:1.25.4 as builder

WORKDIR /app

COPY .. .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM scratch

COPY --from=builder server /bin/server
COPY --from=builder /app/configs/config.yml /bin/config.yml
COPY --from=builder /app/docs /docs

ENTRYPOINT ["/bin/server", "--config=/bin/config.yml"]