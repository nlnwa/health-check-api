# Health Checker

## Health check API for Veidemann

## Build

```bash
go build
```

## Run

```bash
./health-check-api
```

The health checker exposes two endpoints both of which responds with [Health Check Response Format for HTTP APIs](https://tools.ietf.org/id/draft-inadarei-api-health-check-03.html)

- **Health endpoint**

    http://localhost:8080/health

- **Liveness endpoint (health of health checker)**
    
    Used by for example kubernetes liveness probe.
    
    http://localhost:8080/healthz


## Configuration

Web endpoints to be checked are configurable and must be specified in a configuration file (defaults to `./config.yaml`).

**config.yaml:**

```yaml
web:
  - name: dashboard
    url: http://host/path
  - name: rest-api
  - url: http://host/api-path
```

Web endpoints are checked by a HTTP HEAD request and considered healthy if response is HTTP status code between 2xx-3xx.
