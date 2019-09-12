package check

import (
	"net/http"

	"github.com/nlnwa/health-check-api/pkg/health"
)

type web struct {
	name string
	url  string
}

func NewWebChecker(name string, url string) health.Checker {
	return web{
		url:  url,
		name: name,
	}
}

func (web web) Name() string {
	return web.name
}

func (web web) Check() []health.CheckResponse {
	check := health.CheckResponse{}

	check.Time = health.GetCurrentTime()

	resp, err := http.Head(web.url)

	if err != nil {
		check.Status = health.StatusUnhealthy
		check.Output = err.Error()
	} else if resp.StatusCode < 400 {
		check.Status = health.StatusHealthy
	} else if resp.StatusCode >= 500 {
		check.Status = health.StatusUnhealthy
		check.Output = resp.Status
	} else {
		check.Status = health.StatusWarn
		check.Output = resp.Status
	}

	return []health.CheckResponse{check}
}
