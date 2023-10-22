FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as build
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /out/backend -ldflags "-s -w"
COPY data /out/data
COPY internal/platform/migrations /out/internal/platform/migrations

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=build /out .
ENTRYPOINT ["./backend", "-bind=:5001", "-callback-addr=http://backend:5001", "-moneropay=http://moneropay:5000", "-postgres=postgresql://postgres:sup3rs3cure@postgresql:5432/metronero?sslmode=disable", "-token-secret=xufeg6uoth4oM7CohK2D", "-callback-url=http://backend:5001/callback"]
