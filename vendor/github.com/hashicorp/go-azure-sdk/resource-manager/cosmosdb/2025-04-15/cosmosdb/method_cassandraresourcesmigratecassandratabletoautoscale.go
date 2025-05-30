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

type CassandraResourcesMigrateCassandraTableToAutoscaleOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ThroughputSettingsGetResults
}

// CassandraResourcesMigrateCassandraTableToAutoscale ...
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToAutoscale(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesMigrateCassandraTableToAutoscaleOperationResponse, err error) {
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

// CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll performs CassandraResourcesMigrateCassandraTableToAutoscale then polls until it's completed
func (c CosmosDBClient) CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll(ctx context.Context, id CassandraKeyspaceTableId) error {
	result, err := c.CassandraResourcesMigrateCassandraTableToAutoscale(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesMigrateCassandraTableToAutoscale: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CassandraResourcesMigrateCassandraTableToAutoscale: %+v", err)
	}

	return nil
}
