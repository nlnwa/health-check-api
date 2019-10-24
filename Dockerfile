FROM docker.io/golang:latest as golang

WORKDIR /build

COPY . .

RUN mkdir -p /out

# Cache builds without version info
RUN go build -mod readonly -o /out/veidemann-health-check-api -ldflags "-s -w"

ARG VERSION
ENV GO_LDFLAGS="-s -w -X github.com/nlnwa/veidemann-health-check-api/pkg/version.Version=${VERSION}"
RUN go build -mod readonly -o /out/veidemann-health-check-api -ldflags "${GO_LDFLAGS}"


FROM gcr.io/distroless/base

COPY LICENSE /LICENSE

COPY --from=golang /out /out

ENTRYPOINT ["/out/veidemann-health-check-api"]

EXPOSE 8080
