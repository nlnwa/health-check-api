FROM golang:1.15 as build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=dirty

# -trimpath remove file system paths from executable
# -ldflags arguments passed to go tool link:
#   -s disable symbol table
#   -w disable DWARF generation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.Version=$VERSION" .


FROM gcr.io/distroless/base

COPY LICENSE /LICENSE

COPY --from=build build/veidemann-health-check-api /

ENTRYPOINT ["/veidemann-health-check-api"]

EXPOSE 8080
