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

type SqlResourcesMigrateSqlDatabaseToAutoscaleOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ThroughputSettingsGetResults
}

// SqlResourcesMigrateSqlDatabaseToAutoscale ...
func (c CosmosDBClient) SqlResourcesMigrateSqlDatabaseToAutoscale(ctx context.Context, id SqlDatabaseId) (result SqlResourcesMigrateSqlDatabaseToAutoscaleOperationResponse, err error) {
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

// SqlResourcesMigrateSqlDatabaseToAutoscaleThenPoll performs SqlResourcesMigrateSqlDatabaseToAutoscale then polls until it's completed
func (c CosmosDBClient) SqlResourcesMigrateSqlDatabaseToAutoscaleThenPoll(ctx context.Context, id SqlDatabaseId) error {
	result, err := c.SqlResourcesMigrateSqlDatabaseToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesMigrateSqlDatabaseToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SqlResourcesMigrateSqlDatabaseToAutoscale: %+v", err)
	}

	return nil
}
