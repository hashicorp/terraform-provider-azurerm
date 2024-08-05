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

type CassandraResourcesCreateUpdateCassandraTableOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CassandraTableGetResults
}

// CassandraResourcesCreateUpdateCassandraTable ...
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraTable(ctx context.Context, id CassandraKeyspaceTableId, input CassandraTableCreateUpdateParameters) (result CassandraResourcesCreateUpdateCassandraTableOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// CassandraResourcesCreateUpdateCassandraTableThenPoll performs CassandraResourcesCreateUpdateCassandraTable then polls until it's completed
func (c CosmosDBClient) CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx context.Context, id CassandraKeyspaceTableId, input CassandraTableCreateUpdateParameters) error {
	result, err := c.CassandraResourcesCreateUpdateCassandraTable(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesCreateUpdateCassandraTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CassandraResourcesCreateUpdateCassandraTable: %+v", err)
	}

	return nil
}
