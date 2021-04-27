package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nlnwa/veidemann-health-check-api/pkg/api"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/controller"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/prometheus"
	"github.com/nlnwa/veidemann-health-check-api/pkg/healthcheck"
	"github.com/nlnwa/veidemann-health-check-api/pkg/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Version = "undefined"

var statusToApi = map[healthcheck.Status]api.Status{
	healthcheck.StatusPass:      api.StatusHealthy,
	healthcheck.StatusWarning:   api.StatusWarn,
	healthcheck.StatusFail:      api.StatusUnhealthy,
	healthcheck.StatusUndefined: api.StatusWarn,
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, no-cache, must-revalidate, max-age=3600")
	w.Header().Set("Content-Type", "application/health+json; charset=UTF-8")
	w.Header().Set("Expires", "0")
	w.Header().Set("Vary", "Accept-Encoding")
}

func healthCollector(health *api.Health) healthcheck.CheckObserver {
	health.Checks = make(map[string][]api.Check)
	health.Status = api.StatusHealthy

	return func(healthCheck *healthcheck.CheckResult) {
		var checks []api.Check
		for _, checkResult := range healthCheck.Results {
			check := api.Check{
				Time:              checkResult.Time.Format(time.RFC3339),
				ComponentType:     checkResult.Type,
				ComponentId:       checkResult.Id,
				Status:            statusToApi[checkResult.Status],
				ObservedUnit:      checkResult.Unit,
				ObservedValue:     checkResult.Value,
				AffectedEndpoints: checkResult.Endpoints,
				Links:             checkResult.Links,
				Output: func(err error) string {
					if err != nil {
						return err.Error()
					}
					return ""
				}(checkResult.Err),
			}
			checks = append(checks, check)
			if check.Status.Value() > 0 && check.Status.Value() < health.Status.Value() {
				health.Status = check.Status
			}
		}
		health.Checks[healthCheck.Name] = checks
	}
}

func healthCheckHandler(hc healthcheck.HealthChecker) http.HandlerFunc {
	health := &api.Health{
		Version:   api.Version,
		ReleaseId: Version,
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		setDefaultHeaders(w)

		hc.RunChecks(healthCollector(health))

		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error().Err(err).Msg("Failed to encode Health struct")
		}
	}
}

// Liveness probe endpoint for the health check API itself
func livenessHandler() http.HandlerFunc {
	response, err := json.Marshal(api.Health{Status: api.StatusHealthy})
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		setDefaultHeaders(w)
		if _, err := w.Write(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type Config struct {
	LogLevel         string `mapstructure:"log-level"`
	LogFormat        string `mapstructure:"log-format"`
	LogMethod        bool   `mapstructure:"log-method"`
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	HealthPath       string `mapstructure:"health-path"`
	LivenessPath     string `mapstructure:"liveness-path"`
	ControllerHost   string `mapstructure:"controller-host"`
	ControllerPort   int    `mapstructure:"controller-port"`
	ControllerApiKey string `mapstructure:"controller-api-key"`
	PrometheusUrl    string `mapstructure:"prometheus-url"`
	VersionsPath     string `mapstructure:"versions-path"`
}

func main() {
	pflag.String("host", "", "Listening interface")
	pflag.Int("port", 8080, "Listening port")
	pflag.String("health-path", "/health", "URL path of health endpoint")
	pflag.String("liveness-path", "/healthz", "URL path of liveness endpoint")
	pflag.String("controller-host", "veidemann-controller", "Veidemann controller host")
	pflag.Int("controller-port", 7700, "Veidemann controller port")
	pflag.String("controller-api-key", "", "Veidemann controller API key")
	pflag.String("prometheus-url", "http://localhost:9090", "Prometheus HTTP API URL")
	pflag.String("config-file", "config", "Name of config file (without extension)")
	pflag.String("config-path", ".", "Path to look for config file in")
	pflag.String("versions-path", "./versions.json", "Path to versions file")
	pflag.String("log-level", "info", "Log level, available levels are: panic, fatal, error, warn, info, debug and trace")
	pflag.String("log-format", "logfmt", "Log format, available values are: logfmt and json")
	pflag.Bool("log-method", false, "Log file:line of method caller")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile(viper.GetString("config-file"))
	viper.AddConfigPath(viper.GetString("config-path"))
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(err)
		}
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	logger.InitLog(config.LogLevel, config.LogFormat, config.LogMethod)

	controllerClient := controller.New(controller.Options{
		Host:   config.ControllerHost,
		Port:   config.ControllerPort,
		ApiKey: config.ControllerApiKey,
	})
	if err := controllerClient.Connect(); err != nil {
		panic(err)
	}

	prometheusClient := prometheus.New(config.PrometheusUrl)

	healthChecker := healthcheck.New(&healthcheck.Options{
		ControllerClient: controllerClient,
		PrometheusClient: prometheusClient,
		VersionPath:      config.VersionsPath,
	})

	router := http.NewServeMux()
	router.HandleFunc(config.LivenessPath, livenessHandler())
	router.HandleFunc(config.HealthPath, healthCheckHandler(healthChecker))

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// wait for signal
		sig := <-signals
		log.Info().Str("signal", sig.String()).Msg("Shutting down")

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("")
		}
	}()

	log.Info().
		Str("version", Version).
		Str("address", srv.Addr).
		Msg("Server listening")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Server failure")
	}
	<-done
}
