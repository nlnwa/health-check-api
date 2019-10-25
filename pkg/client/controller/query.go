package controller

import (
	"context"

	"github.com/nlnwa/veidemann-api-go/frontier/v1"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"github.com/pkg/errors"
)

type Query interface {
	GetRunningJobs(ctx context.Context) ([]string, error)
}

// RunLanguageDetection calls the gRPC method with the same name.
func (ac Client) listJobExecutions(ctx context.Context) (*api.JobExecutionsListReply, error) {
	conn, err := ac.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := api.NewStatusClient(conn)

	req := &api.ListJobExecutionsRequest{}

	reply, err := client.ListJobExecutions(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list job executions")
	}
	return reply, nil
}

func (ac Client) GetRunningJobs(ctx context.Context) ([]string, error) {
	var ids []string

	reply, err := ac.listJobExecutions(ctx)
	if err != nil {
		return ids, err
	}
	for _, value := range reply.GetValue() {
		if value.GetState() == frontier.JobExecutionStatus_RUNNING {
			ids = append(ids, value.Id)
		}
	}
	return ids, nil
}
