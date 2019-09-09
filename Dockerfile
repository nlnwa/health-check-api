FROM golang:alpine as golang

WORKDIR /build

COPY . .

RUN mkdir -p /out

# Cache builds without version info
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/health-check-api -ldflags "-s -w"

ARG VERSION
ENV GO_LDFLAGS="-s -w -X github.com/nlnwa/nettarkivet-health-check-api/pkg/version.Version=${VERSION}"
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/health-check-api -ldflags "${GO_LDFLAGS}"


FROM gcr.io/distroless/base

COPY LICENSE /LICENSE

COPY --from=golang /out /out

ENTRYPOINT ["/out/health-check-api"]

EXPOSE 8080
