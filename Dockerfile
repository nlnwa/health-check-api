FROM docker.io/golang:1.14 as golang

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Cache build without version info
RUN go build -trimpath -ldflags "-s -w"

ARG VERSION
ENV GO_LDFLAGS="-s -w -X github.com/nlnwa/veidemann-health-check-api/pkg/version.Version=${VERSION}"
RUN go build -trimpath -ldflags "${GO_LDFLAGS}"


FROM gcr.io/distroless/base

COPY LICENSE /LICENSE

COPY --from=golang /build/veidemann-health-check-api /

ENTRYPOINT ["/veidemann-health-check-api"]

EXPOSE 8080
