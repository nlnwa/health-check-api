package health

// https://tools.ietf.org/html/draft-inadarei-api-health-check-03#section-3

type status int

const (
	statusUnhealthy status = iota
	statusWarn
	statusHealthy
)

var statusMapStandard = map[status]string{
	statusUnhealthy: "fail",
	statusWarn:      "warn",
	statusHealthy:   "pass",
}

var statusMap = statusMapStandard

var StatusHealthy = statusMap[statusHealthy]

var StatusUnhealthy = statusMap[statusUnhealthy]

var StatusWarn = statusMap[statusWarn]

type Response struct {
	Status      string                     `json:"status"` // mandatory
	Version     string                     `json:"version,omitempty"`
	ReleaseId   string                     `json:"releaseId,omitempty"`
	Notes       string                     `json:"notes,omitempty"`
	Output      string                     `json:"output,omitempty"`
	Checks      map[string][]CheckResponse `json:"checks,omitempty"`
	Links       []string                   `json:"links,omitempty"`
	ServiceId   string                     `json:"serviceId,omitempty"`
	Description string                     `json:"description,omitempty"`
}

type CheckResponse struct {
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

type Checker interface {
	Check() []CheckResponse
	Name() string
}
