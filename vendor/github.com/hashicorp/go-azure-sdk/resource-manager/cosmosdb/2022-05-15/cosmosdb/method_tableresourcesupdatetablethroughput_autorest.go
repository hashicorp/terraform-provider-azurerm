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

type TableResourcesUpdateTableThroughputOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TableResourcesUpdateTableThroughput ...
func (c CosmosDBClient) TableResourcesUpdateTableThroughput(ctx context.Context, id TableId, input ThroughputSettingsUpdateParameters) (result TableResourcesUpdateTableThroughputOperationResponse, err error) {
	req, err := c.preparerForTableResourcesUpdateTableThroughput(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesUpdateTableThroughput", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTableResourcesUpdateTableThroughput(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "TableResourcesUpdateTableThroughput", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TableResourcesUpdateTableThroughputThenPoll performs TableResourcesUpdateTableThroughput then polls until it's completed
func (c CosmosDBClient) TableResourcesUpdateTableThroughputThenPoll(ctx context.Context, id TableId, input ThroughputSettingsUpdateParameters) error {
	result, err := c.TableResourcesUpdateTableThroughput(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TableResourcesUpdateTableThroughput: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TableResourcesUpdateTableThroughput: %+v", err)
	}

	return nil
}

// preparerForTableResourcesUpdateTableThroughput prepares the TableResourcesUpdateTableThroughput request.
func (c CosmosDBClient) preparerForTableResourcesUpdateTableThroughput(ctx context.Context, id TableId, input ThroughputSettingsUpdateParameters) (*http.Request, error) {
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

// senderForTableResourcesUpdateTableThroughput sends the TableResourcesUpdateTableThroughput request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForTableResourcesUpdateTableThroughput(ctx context.Context, req *http.Request) (future TableResourcesUpdateTableThroughputOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
