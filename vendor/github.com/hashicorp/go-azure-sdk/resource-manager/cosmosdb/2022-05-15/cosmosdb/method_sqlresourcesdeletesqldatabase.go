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

type SqlResourcesDeleteSqlDatabaseOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// SqlResourcesDeleteSqlDatabase ...
func (c CosmosDBClient) SqlResourcesDeleteSqlDatabase(ctx context.Context, id SqlDatabaseId) (result SqlResourcesDeleteSqlDatabaseOperationResponse, err error) {
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

// SqlResourcesDeleteSqlDatabaseThenPoll performs SqlResourcesDeleteSqlDatabase then polls until it's completed
func (c CosmosDBClient) SqlResourcesDeleteSqlDatabaseThenPoll(ctx context.Context, id SqlDatabaseId) error {
	result, err := c.SqlResourcesDeleteSqlDatabase(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesDeleteSqlDatabase: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after SqlResourcesDeleteSqlDatabase: %+v", err)
	}

	return nil
}
