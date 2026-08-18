// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &mongoClusterEarliestRestoreTimePoller{}

type mongoClusterEarliestRestoreTimePoller struct {
	client *mongoclusters.MongoClustersClient
	id     mongoclusters.MongoClusterId
}

func NewMongoClusterEarliestRestoreTimePoller(client *mongoclusters.MongoClustersClient, id mongoclusters.MongoClusterId) pollers.PollerType {
	return &mongoClusterEarliestRestoreTimePoller{
		client: client,
		id:     id,
	}
}

func (p mongoClusterEarliestRestoreTimePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s while waiting for `earliest_restore_time`: %+v", p.id, err)
	}

	status := pollers.PollingStatusInProgress
	if model := resp.Model; model != nil && model.Properties != nil {
		if backup := model.Properties.Backup; backup != nil && pointer.From(backup.EarliestRestoreTime) != "" {
			status = pollers.PollingStatusSucceeded
		}
	}

	return &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       status,
	}, nil
}
