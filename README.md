# Veidemann Health Checker

## Health check API for Veidemann


The health checker exposes two endpoints both of which responds with [Health Check Response Format for HTTP APIs](https://tools.ietf.org/id/draft-inadarei-api-health-check-03.html):

1. **Health endpoint**

2. **Liveness endpoint (liveness of health checker)**

## Build

```bash
go build
```

## Run

To see what options are available run:

```bash
./veidemann-health-check-api --help
```

## Configuration

Options can be configured via:

1. Configuration file

    Supported formats are JSON, TOML, YAML, HCL, envfile (.env) and Java properties (.conf).
    
    Example:
    
    ```angular2
    # config.yaml
 
    controller-address: "myhost:7700"
    ```

2. Environment variables

    Environment variables take precedence over values in configuration file.
    
    ```bash
    CONTROLLER_ADDRESS=myhost:7700 ./veidemann-health-check-api
    ```

3. Flags

    Flags take precedence over environment variables.
    
    ```bash
    ./veidemann-health-check-api --controller-api-key ABCD-1234
    ```


## Skaffold

The `k8s` folder contains kubernetes manifests used by the _skaffold_
configuration `skaffold.yaml`.

```bash
skaffold dev --port-forward
```
