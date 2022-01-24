FROM golang as builder
WORKDIR /workspace/app
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM gcr.io/distroless/static:nonroot as final
WORKDIR /
COPY --from=builder /workspace/app .
USER 65532:65532
ENTRYPOINT ["/app"]
