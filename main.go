package main

import (
	"context"
	"encoding/json"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nlnwa/veidemann-health-check-api/pkg/api"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/controller"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/prometheus"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/web"
	"github.com/nlnwa/veidemann-health-check-api/pkg/healthcheck"
	"github.com/nlnwa/veidemann-health-check-api/pkg/version"
	"github.com/spf13/viper"
)

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, no-cache, must-revalidate, max-age=3600")
	w.Header().Set("Content-Type", "application/health+json; charset=UTF-8")
	w.Header().Set("Expires", "0")
	w.Header().Set("Vary", "Accept-Encoding")
}

func healthCollector(health *api.Health) func(*healthcheck.CheckResult) {
	health.Status = api.StatusHealthy
	health.Version = version.Version
	health.Checks = make(map[string][]api.Check, 0)

	return func(healthCheck *healthcheck.CheckResult) {
		var checks []api.Check
		for _, checkResult := range healthCheck.Results {
			check := api.Check{
				Time:              api.GetCurrentTime(),
				ComponentType:     checkResult.Type,
				ComponentId:       checkResult.Id,
				Status:            statusToApi(checkResult.Status),
				ObservedUnit:      checkResult.Unit,
				ObservedValue:     checkResult.Value,
				AffectedEndpoints: checkResult.Endpoints,
				Links:             checkResult.Links,
				Description:       checkResult.Description,
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

func statusToApi(status healthcheck.Status) api.Status {
	statusToApi := map[healthcheck.Status]api.Status{
		healthcheck.StatusPass:    api.StatusHealthy,
		healthcheck.StatusWarning: api.StatusWarn,
		healthcheck.StatusFail:    api.StatusUnhealthy,
	}
	return statusToApi[status]
}

func healthCheckHandler(hc *healthcheck.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		setDefaultHeaders(w)

		health := &api.Health{}

		hc.RunChecks(healthCollector(health))

		writer := io.MultiWriter(w, log.Writer())

		if err := json.NewEncoder(writer).Encode(health); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
	}
}

// Liveness probe endpoint for the health check API itself
func livenessHandler() http.HandlerFunc {
	response, err := json.Marshal(api.Health{Status: api.StatusHealthy})
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		setDefaultHeaders(w)
		if _, err := w.Write(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type Config struct {
	Port                  string
	HealthPath            string `mapstructure:"health-path"`
	LivenessPath          string `mapstructure:"liveness-path"`
	VeidemannDashboardUrl string `mapstructure:"veidemann-dashboard-url"`
	ControllerHost        string `mapstructure:"controller-host"`
	ControllerPort        int    `mapstructure:"controller-port"`
	ControllerApiKey      string `mapstructure:"controller-api-key"`
	PrometheusUrl         string `mapstructure:"prometheus-url"`
}

func main() {
	// configuration defaults
	port := "8080"
	healthPath := "/health"
	livenessPath := "/healthz"
	configFileName := "config"
	configPath := "."
	controllerHost := "veidemann-controller"
	controllerPort := 7700
	controllerApiKey := ""
	prometheusUrl := "http://localhost:9090"
	veidemannDashboardUrl := "http://localhost/veidemann"

	flag.StringVar(&port, "port", port, "Listening port")
	flag.StringVar(&healthPath, "health-path", healthPath, "URL path of health endpoint")
	flag.StringVar(&livenessPath, "liveness-path", livenessPath, "URL path of liveness endpoint")
	flag.StringVar(&veidemannDashboardUrl, "veidemann-dashboard-url", veidemannDashboardUrl, "URL of veidemann dashboard")
	flag.StringVar(&controllerHost, "controller-host", controllerHost, "Veidemann controller host")
	flag.IntVar(&controllerPort, "controller-port", controllerPort, "Veidemann controller port")
	flag.StringVar(&controllerApiKey, "controller-api-key", controllerApiKey, "Veidemann controller API key")
	flag.StringVar(&prometheusUrl, "prometheus-url", prometheusUrl, "Prometheus HTTP API URL")
	flag.StringVar(&configFileName, "config-file", configFileName, "Name of config file (without extension)")
	flag.StringVar(&configPath, "config-path", configPath, "Path to look for config file in")
	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile(configFileName)
	viper.AddConfigPath(configPath)
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

	healthCheckerOptions := healthcheck.Options{
		Controller: controller.Options{
			Host:   config.ControllerHost,
			Port:   config.ControllerPort,
			ApiKey: config.ControllerApiKey,
		},
		WebOptions: web.Options{
			VeidemannDashboardUrl: config.VeidemannDashboardUrl,
		},
		Prometheus: prometheus.Options{
			Address: config.PrometheusUrl,
		},
	}
	healthChecker := healthcheck.NewHealthChecker(&healthCheckerOptions)

	router := http.NewServeMux()
	router.HandleFunc(config.LivenessPath, livenessHandler())
	router.HandleFunc(config.HealthPath, healthCheckHandler(healthChecker))

	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// shutdown gracefully
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// wait for signal
		<-done

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
