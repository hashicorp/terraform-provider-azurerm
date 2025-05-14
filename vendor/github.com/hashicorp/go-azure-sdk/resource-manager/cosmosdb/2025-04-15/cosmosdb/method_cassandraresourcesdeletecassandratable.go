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

type CassandraResourcesDeleteCassandraTableOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// CassandraResourcesDeleteCassandraTable ...
func (c CosmosDBClient) CassandraResourcesDeleteCassandraTable(ctx context.Context, id CassandraKeyspaceTableId) (result CassandraResourcesDeleteCassandraTableOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodDelete,
		Path:       id.ID(),
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

// CassandraResourcesDeleteCassandraTableThenPoll performs CassandraResourcesDeleteCassandraTable then polls until it's completed
func (c CosmosDBClient) CassandraResourcesDeleteCassandraTableThenPoll(ctx context.Context, id CassandraKeyspaceTableId) error {
	result, err := c.CassandraResourcesDeleteCassandraTable(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraResourcesDeleteCassandraTable: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CassandraResourcesDeleteCassandraTable: %+v", err)
	}

	return nil
}
