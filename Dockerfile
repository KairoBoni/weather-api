FROM golang:1.22 as builder

#avoid root
ENV USER=appuser
ENV UID=1000

#avoid root
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GO111MODULES=on go build -o weather-api ./app/weather/

FROM scratch

#avoid rootless
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app
COPY --from=builder /app/weather-api /app

#avoid rootless
USER appuser:appuser

ENTRYPOINT ["/app/weather-api"]