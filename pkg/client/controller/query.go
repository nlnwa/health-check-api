package controller

import (
	"context"
	"fmt"
	"github.com/nlnwa/veidemann-api/go/commons/v1"
	"github.com/nlnwa/veidemann-api/go/config/v1"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemann-api/go/controller/v1"
	"github.com/nlnwa/veidemann-api/go/frontier/v1"
	"github.com/nlnwa/veidemann-api/go/report/v1"
)

type Query interface {
	ListFetchingSeeds(ctx context.Context, pageSize int32) ([]string, error)
	GetCrawlerStatus(ctx context.Context) (*controller.CrawlerStatus, error)
}

func (ac Client) ListFetchingSeeds(ctx context.Context, pageSize int32) ([]string, error) {
	reportClient := report.NewReportClient(ac.conn)
	configClient := config.NewConfigClient(ac.conn)

	req := &report.CrawlExecutionsListRequest{
		PageSize:           pageSize,
		ReturnedFieldsMask: &commons.FieldMask{Paths: []string{"seedId"}},
		State: []frontier.CrawlExecutionStatus_State{
			frontier.CrawlExecutionStatus_FETCHING,
		},
	}
	stream, err := reportClient.ListExecutions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list crawl executions: %w", err)
	}
	var statuses []string
	for {
		status, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status.GetSeedId())
	}
	if len(statuses) == 0 {
		return statuses, nil
	}
	seedReq := &config.ListRequest{
		Id:                 statuses,
		Kind:               config.Kind_seed,
		ReturnedFieldsMask: &commons.FieldMask{Paths: []string{"meta.name"}},
	}

	co, err := configClient.ListConfigObjects(ctx, seedReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list seeds: %w", err)
	}
	var seeds []string
	for {
		seed, err := co.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		seeds = append(seeds, seed.Meta.Name)
	}
	return seeds, nil
}

func (ac Client) GetCrawlerStatus(ctx context.Context) (*controller.CrawlerStatus, error) {
	client := controller.NewControllerClient(ac.conn)

	status, err := client.Status(ctx, &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to check status: %w", err)
	}
	return status, nil
}
