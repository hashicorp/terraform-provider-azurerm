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

type MongoDBResourcesMigrateMongoDBCollectionToManualThroughputOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ThroughputSettingsGetResults
}

// MongoDBResourcesMigrateMongoDBCollectionToManualThroughput ...
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBCollectionToManualThroughput(ctx context.Context, id MongodbDatabaseCollectionId) (result MongoDBResourcesMigrateMongoDBCollectionToManualThroughputOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/throughputSettings/default/migrateToManualThroughput", id.ID()),
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

// MongoDBResourcesMigrateMongoDBCollectionToManualThroughputThenPoll performs MongoDBResourcesMigrateMongoDBCollectionToManualThroughput then polls until it's completed
func (c CosmosDBClient) MongoDBResourcesMigrateMongoDBCollectionToManualThroughputThenPoll(ctx context.Context, id MongodbDatabaseCollectionId) error {
	result, err := c.MongoDBResourcesMigrateMongoDBCollectionToManualThroughput(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesMigrateMongoDBCollectionToManualThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesMigrateMongoDBCollectionToManualThroughput: %+v", err)
	}

	return nil
}
