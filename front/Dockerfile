FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

# Static binary: CGO disabled so it runs on distroless "static" (no libc).
# Assets (templates/static) are embedded via go:embed in app/web, so the
# runtime image needs nothing besides this binary.
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
