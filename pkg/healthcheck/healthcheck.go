package healthcheck

import (
	"context"
	"fmt"
	"github.com/nlnwa/veidemann-health-check-api/pkg/version"
	"log"
	"time"

	controllerApi "github.com/nlnwa/veidemann-api/go/controller/v1"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/controller"
	"github.com/nlnwa/veidemann-health-check-api/pkg/client/prometheus"
)

const (
	VeidemannFetching string = "veidemann:fetching"
	VeidemannStatus   string = "veidemann:status"
	VeidemannActivity string = "veidemann:activity"
	VeidemannHarvest  string = "veidemann:harvest"
	VeidemannVersion  string = "veidemann:version"
)

type Value interface {
}

type Status int

const (
	StatusUndefined = iota
	StatusPass
	StatusWarning
	StatusFail
)

type CheckResult struct {
	Name    string
	Results []*Result
}

type Result struct {
	Id        string
	Type      string
	Value     Value
	Unit      string
	Status    Status
	Endpoints []string
	Links     []string
	Time      time.Time
	Err       error
}

type checker func(context.Context) *Result

type component struct {
	id       string
	checkers []checker
}

type CheckObserver func(*CheckResult)

type Options struct {
	Controller  controller.Options
	Prometheus  prometheus.Options
	VersionPath string
}

type HealthCheckere interface {
	RunChecks(CheckObserver)
}

type HealthChecker struct {
	prometheusClient prometheus.Query
	controllerClient controller.Query
	versionPath      string
	components       []component
}

func New(options *Options) HealthCheckere {
	hc := &HealthChecker{
		controllerClient: controller.New(options.Controller),
		prometheusClient: prometheus.New(options.Prometheus),
		versionPath:      options.VersionPath,
	}
	hc.components = hc.getChecks()
	return hc
}

func (hc *HealthChecker) RunChecks(observer CheckObserver) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for _, component := range hc.components {
		var checkResults []*Result

		for _, checker := range component.checkers {
			checkResults = append(checkResults, checker(ctx))
		}
		result := &CheckResult{
			Name:    component.id,
			Results: checkResults,
		}
		observer(result)
	}
}

// getChecks returns a list of components to be checked
func (hc *HealthChecker) getChecks() []component {
	var runStatus controllerApi.RunStatus
	var isActive bool
	var nrOfFetchingSeeds int
	var queueSize int64
	var versions *Result
	return []component{
		{
			id: VeidemannVersion,
			checkers: []checker{
				func(ctx context.Context) *Result {
					if versions != nil {
						return versions
					}
					value, err := version.GetVersions(hc.versionPath)
					versions = &Result{
						Time:  time.Now(),
						Err:   err,
						Value: value,
						Status: func(err error) Status {
							if err != nil {
								return StatusWarning
							}
							return StatusPass
						}(err),
					}
					return versions
				},
			},
		},
		{
			id: VeidemannStatus,
			checkers: []checker{
				func(ctx context.Context) *Result {
					crawlerStatus, err := hc.controllerClient.GetCrawlerStatus(ctx)
					if err != nil {
						log.Println(err)
					}
					runStatus = crawlerStatus.GetRunStatus()
					queueSize = crawlerStatus.GetQueueSize()
					result := &Result{
						Type:  "harvester",
						Time:  time.Now(),
						Err:   err,
						Value: fmt.Sprintf("%v", crawlerStatus),
						Status: func(err error) Status {
							if err != nil {
								return StatusUndefined
							}
							return StatusPass
						}(err),
					}
					return result
				},
			},
		},
		{
			id: VeidemannFetching,
			checkers: []checker{
				func(ctx context.Context) *Result {
					fetchingSeeds, err := hc.controllerClient.ListFetchingSeeds(ctx, 5)
					if err != nil {
						log.Println(err)
					}
					nrOfFetchingSeeds = len(fetchingSeeds)
					result := &Result{
						Type:  "harvester",
						Time:  time.Now(),
						Err:   err,
						Value: fetchingSeeds,
						Status: func(err error) Status {
							if err != nil {
								return StatusUndefined
							}
							return StatusPass
						}(err),
					}
					return result
				},
			},
		},
		{
			id: VeidemannActivity,
			checkers: []checker{
				func(ctx context.Context) *Result {
					isActivity, err := hc.prometheusClient.IsActivity(ctx)
					result := &Result{
						Type:  "harvester",
						Time:  time.Now(),
						Err:   err,
						Value: isActivity,
						Status: func(err error) Status {
							if err != nil {
								return StatusUndefined
							}
							return StatusPass
						}(err),
					}
					return result
				},
			},
		},
		{
			id: VeidemannHarvest,
			checkers: []checker{
				func(ctx context.Context) *Result {
					return &Result{
						Type: "harvester",
						Time: time.Now(),
						Status: func() Status {
							switch runStatus {
							case controllerApi.RunStatus_RUNNING:
								if queueSize > 0 && isActive && nrOfFetchingSeeds == 0 {
									return StatusFail
								}
							}
							return StatusPass
						}(),
					}
				},
			},
		},
	}
}
