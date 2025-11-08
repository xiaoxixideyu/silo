FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /out/admin ./cmd/admin/main.go
RUN go build -o /out/website ./cmd/website/main.go

FROM alpine:latest
ARG QQFUN_APP_SERVER_TYPE=install
ENV QQFUN_APP_SERVER_TYPE=${QQFUN_APP_SERVER_TYPE}
RUN apk --no-cache add ca-certificates
COPY --from=builder /out/admin /app/admin
COPY --from=builder /out/website /app/website
WORKDIR /app
ENTRYPOINT ["./${QQFUN_APP_SERVER_TYPE}"]
CMD ["start"]