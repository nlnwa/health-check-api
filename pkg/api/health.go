package api

// https://tools.ietf.org/html/draft-inadarei-api-health-check-03#section-3

type StatusCode int

type Status string

type Value interface{}
type Checks map[string][]Check

const (
	statusUndefined StatusCode = iota
	statusUnhealthy
	statusWarn
	statusHealthy
)

var statusMap = map[StatusCode]Status{
	statusUndefined: "",
	statusUnhealthy: "down",
	statusWarn:      "warn",
	statusHealthy:   "up",
}

var statusToStatusCode = map[Status]StatusCode{
	"down": statusUnhealthy,
	"warn": statusWarn,
	"up":   statusHealthy,
	"":     statusUndefined,
}

var StatusHealthy = statusMap[statusHealthy]
var StatusUnhealthy = statusMap[statusUnhealthy]
var StatusWarn = statusMap[statusWarn]

func (s Status) Value() StatusCode {
	return statusToStatusCode[s]
}

type Check struct {
	ComponentId       string   `json:"componentId,omitempty"`
	ComponentType     string   `json:"componentType,omitempty"`
	ObservedValue     Value    `json:"observedValue,omitempty"`
	ObservedUnit      string   `json:"observedUnit,omitempty"`
	Status            Status   `json:"status,omitempty"`
	AffectedEndpoints []string `json:"affectedEndpoints,omitempty"`
	Time              string   `json:"time,omitempty"`
	Output            string   `json:"output,omitempty"`
	Links             []string `json:"links,omitempty"`
	Description       string   `json:"description,omitempty"`
}

type Health struct {
	Status      Status   `json:"status"` // mandatory
	Version     string   `json:"version,omitempty"`
	ReleaseId   string   `json:"releaseId,omitempty"`
	Notes       []string   `json:"notes,omitempty"`
	Output      string   `json:"output,omitempty"`
	Checks      Checks   `json:"checks,omitempty"`
	Links       []string `json:"links,omitempty"`
	ServiceId   string   `json:"serviceId,omitempty"`
	Description string   `json:"description,omitempty"`
}
