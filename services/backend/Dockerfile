# Multi-stage build - static binary in a distroless runtime, no shell/package manager in the final image.
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/backend .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/backend /backend
EXPOSE 8080
ENTRYPOINT ["/backend"]
