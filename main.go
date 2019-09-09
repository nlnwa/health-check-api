package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// https://tools.ietf.org/html/draft-inadarei-api-health-check-03#section-3

type Status int

const (
	StatusFail Status = iota
	StatusWarn
	StatusPass
)

//var statusMapStandard = map[Status]string{
//	StatusFail: "unhealthy",
//	StatusWarn: "warn",
//	StatusPass: "healthy",
//}

var statusMapSpringBoot = map[Status]string{
	StatusFail: "down",
	StatusPass: "up",
	StatusWarn: "warn",
}

var statusMap = statusMapSpringBoot

func Pass() string {
	return statusMap[StatusPass]
}

func Fail() string {
	return statusMap[StatusFail]
}

func Warn() string {
	return statusMap[StatusWarn]
}

type Health struct {
	Status      string           `json:"status"` // single mandatory root field
	Version     string           `json:"version,omitempty"`
	ReleaseId   string           `json:"releaseId,omitempty"`
	Notes       string           `json:"notes,omitempty"`
	Output      string           `json:"output,omitempty"`
	Checks      map[string]Check `json:"checks,omitempty"`
	Links       []string         `json:"links,omitempty"`
	ServiceId   string           `json:"serviceId,omitempty"`
	Description string           `json:"description,omitempty"`
}

type Check struct {
	ComponentId       string   `json:"componentId,omitempty"`
	ComponentType     string   `json:"componentType,omitempty"`
	ObservedValue     string   `json:"observedValue,omitempty"`
	ObservedUnit      string   `json:"observedUnit,omitempty"`
	Status            string   `json:"status,omitempty"`
	AffectedEndpoints []string `json:"affectedEndpoints,omitempty"`
	Time              string   `json:"time,omitempty"`
	Output            string   `json:"output,omitempty"`
	Links             []string `json:"links,omitempty"`
}

func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Expires", "0")
	w.Header().Set("Vary", "Accept-Encoding")
}

func checkDashboard() (check Check) {
	resp, err := http.Head("https://nettarkivet.nb.no/veidemann")

	check.ComponentType = "component"
	check.Time = getCurrentTime()

	if err != nil {
		check.Status = Fail()
		check.Output = err.Error()
	} else if resp.StatusCode < 400 {
		check.Status = Pass()
	} else {
		check.Status = Warn()
		check.Output = resp.Status
	}

	return
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	setDefaultHeaders(w)

	checks := map[string]Check{
		"veidemann-dashboard": checkDashboard(),
	}

	health := Health{
		Status: statusMap[StatusPass],
		Checks: checks,
	}

	for _, check := range checks {
		if check.Status < health.Status {
			health.Status = check.Status
		}
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", StatusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
