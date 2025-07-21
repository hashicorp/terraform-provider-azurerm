// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &databaseDeletePoller{}

type DatabaseClientInterface interface {
	Get(ctx context.Context, id databases.DatabaseId) (databases.GetOperationResponse, error)
}

type databaseDeletePoller struct {
	databaseClient DatabaseClientInterface
	clusterClient  RedisEnterpriseClientInterface
	databaseId     databases.DatabaseId
	clusterId      redisenterprise.RedisEnterpriseId
}

var (
	deletionSuccess = pollers.PollResult{
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	deletionInProgress = pollers.PollResult{
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewDatabaseDeletePoller(databaseClient DatabaseClientInterface, clusterClient RedisEnterpriseClientInterface, databaseId databases.DatabaseId, clusterId redisenterprise.RedisEnterpriseId) *databaseDeletePoller {
	return &databaseDeletePoller{
		databaseClient: databaseClient,
		clusterClient:  clusterClient,
		databaseId:     databaseId,
		clusterId:      clusterId,
	}
}

func (p databaseDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	log.Printf("[DEBUG] DatabaseDeletePoller checking deletion status for database: %s and cluster: %s", p.databaseId, p.clusterId)

	cluster, err := p.clusterClient.Get(ctx, p.clusterId)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			log.Printf("[DEBUG] DatabaseDeletePoller: cluster %s not found, deletion successful", p.clusterId)
			return &deletionSuccess, nil
		}
		return nil, fmt.Errorf("retrieving cluster %s: %+v", p.clusterId, err)
	}

	db, err := p.databaseClient.Get(ctx, p.databaseId)
	if err != nil {
		if response.WasNotFound(db.HttpResponse) {
			log.Printf("[DEBUG] DatabaseDeletePoller: database %s not found, deletion successful", p.databaseId)
			return &deletionSuccess, nil
		}
		return nil, fmt.Errorf("retrieving database %s: %+v", p.databaseId, err)
	}

	log.Printf("[DEBUG] DatabaseDeletePoller: database %s still exists, deletion in progress", p.databaseId)
	return &deletionInProgress, nil
}
