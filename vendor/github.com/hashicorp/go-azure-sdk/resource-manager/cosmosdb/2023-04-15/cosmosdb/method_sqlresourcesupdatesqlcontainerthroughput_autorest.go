package cosmosdb

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlResourcesUpdateSqlContainerThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesUpdateSqlContainerThroughput ...
func (c CosmosDBClient) SqlResourcesUpdateSqlContainerThroughput(ctx context.Context, id ContainerId, input ThroughputSettingsUpdateParameters) (result SqlResourcesUpdateSqlContainerThroughputOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesUpdateSqlContainerThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesUpdateSqlContainerThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesUpdateSqlContainerThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesUpdateSqlContainerThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesUpdateSqlContainerThroughputThenPoll performs SqlResourcesUpdateSqlContainerThroughput then polls until it's completed
func (c CosmosDBClient) SqlResourcesUpdateSqlContainerThroughputThenPoll(ctx context.Context, id ContainerId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.SqlResourcesUpdateSqlContainerThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesUpdateSqlContainerThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesUpdateSqlContainerThroughput: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesUpdateSqlContainerThroughput prepares the SqlResourcesUpdateSqlContainerThroughput request.
func (c CosmosDBClient) preparerForSqlResourcesUpdateSqlContainerThroughput(ctx context.Context, id ContainerId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/throughputSettings/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSqlResourcesUpdateSqlContainerThroughput sends the SqlResourcesUpdateSqlContainerThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesUpdateSqlContainerThroughput(ctx context.Context, req *http.Request) (future SqlResourcesUpdateSqlContainerThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
