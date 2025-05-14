package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesMigrateGremlinGraphToAutoscaleOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ThroughputSettingsGetResults
}

// GremlinResourcesMigrateGremlinGraphToAutoscale ...
func (c CosmosDBClient) GremlinResourcesMigrateGremlinGraphToAutoscale(ctx context.Context, id GraphId) (result GremlinResourcesMigrateGremlinGraphToAutoscaleOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/throughputSettings/default/migrateToAutoscale", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll performs GremlinResourcesMigrateGremlinGraphToAutoscale then polls until it's completed
func (c CosmosDBClient) GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll(ctx context.Context, id GraphId) error {
	result, err := c.GremlinResourcesMigrateGremlinGraphToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GremlinResourcesMigrateGremlinGraphToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GremlinResourcesMigrateGremlinGraphToAutoscale: %+v", err)
	}

	return nil
}
