# Multi-stage build - static binary in a distroless runtime, no shell/package manager in the final image.
FROM golang:1.22 AS build
WORKDIR /src
COPY . .
# go.sum isn't committed (no local Go toolchain to generate it) - `go mod tidy` resolves and writes it here.
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /out/backend .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/backend /backend
EXPOSE 8080
ENTRYPOINT ["/backend"]
