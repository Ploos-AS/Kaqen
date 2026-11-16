FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -o /kaqen ./cmd/kaqen

FROM alpine:3.21
RUN addgroup -S kaqen && adduser -S -G kaqen kaqen
USER kaqen
COPY --from=build /kaqen /usr/local/bin/kaqen
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/kaqen"]
