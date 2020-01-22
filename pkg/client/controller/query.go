package controller

import (
	"context"
	"fmt"
	"io"

	"github.com/nlnwa/veidemann-api-go/frontier/v1"
	"github.com/nlnwa/veidemann-api-go/report/v1"
)

type Query interface {
	GetRunningJobs(ctx context.Context) ([]string, error)
}

func (ac Client) listRunningJobExecutionStatuses(ctx context.Context) ([]*frontier.JobExecutionStatus, error) {
	conn, err := ac.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	client := report.NewReportClient(conn)

	req := &report.JobExecutionsListRequest{
		State: []frontier.JobExecutionStatus_State{
			frontier.JobExecutionStatus_RUNNING,
		},
	}

	stream, err := client.ListJobExecutions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list job executions: %w", err)
	}
	var jeses []*frontier.JobExecutionStatus
	for {
		jobExecutionStatus, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		jeses = append(jeses, jobExecutionStatus)
	}
	return jeses, nil
}

func (ac Client) GetRunningJobs(ctx context.Context) ([]string, error) {
	var ids []string

	jeses, err := ac.listRunningJobExecutionStatuses(ctx)
	if err != nil {
		return ids, err
	}
	for _, jes := range jeses {
		ids = append(ids, jes.GetId())
	}
	return ids, nil
}
