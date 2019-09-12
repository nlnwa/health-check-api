package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/nlnwa/health-check-api/pkg/check"
	"github.com/nlnwa/health-check-api/pkg/health"
)

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Expires", "0")
	w.Header().Set("Vary", "Accept-Encoding")
}

func healthCheckHandler(healthChecks []health.Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		setDefaultHeaders(w)

		checks := make(map[string][]health.CheckResponse)

		for _, healthCheck := range healthChecks {
			checks[healthCheck.Name()] = healthCheck.Check()
		}

		response := health.Response{
			Status: health.StatusHealthy,
			Checks: checks,
		}

		for _, c := range response.Checks {
			for _, checkResponse := range c {
				if checkResponse.Status < response.Status {
					response.Status = checkResponse.Status
				}
			}
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func livenessHandler() http.HandlerFunc {
	response, err := json.Marshal(health.Response{Status: health.StatusHealthy})
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

type webConfig struct {
	Name string
	Url  string
}

type config struct {
	Web []webConfig
}

func loadConfig(path string) (c config) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	// default values
	port := "8080"
	healthPath := "/health"
	livenessPath := "/healthz"
	configPath := "./config.yaml"

	flag.StringVar(&port, "port", port, "port to listen on")
	flag.StringVar(&healthPath, "check-path", healthPath, "URL path of health endpoint")
	flag.StringVar(&livenessPath, "liveness-path", livenessPath, "URL path of liveness endpoint")
	flag.StringVar(&configPath, "config", configPath, "config file")
	flag.Parse()

	config := loadConfig(configPath)

	var checkers []health.Checker

	for _, value := range config.Web {
		checkers = append(checkers, check.NewWebChecker(value.Name, value.Url))
	}

	http.Handle(livenessPath, livenessHandler())
	http.Handle(healthPath, healthCheckHandler(checkers))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
