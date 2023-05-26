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

type DatabaseAccountsOfflineRegionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabaseAccountsOfflineRegion ...
func (c CosmosDBClient) DatabaseAccountsOfflineRegion(ctx context.Context, id DatabaseAccountId, input RegionForOnlineOffline) (result DatabaseAccountsOfflineRegionOperationResponse, err error) {
	req, err := c.preparerForDatabaseAccountsOfflineRegion(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsOfflineRegion", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabaseAccountsOfflineRegion(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "DatabaseAccountsOfflineRegion", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabaseAccountsOfflineRegionThenPoll performs DatabaseAccountsOfflineRegion then polls until it's completed
func (c CosmosDBClient) DatabaseAccountsOfflineRegionThenPoll(ctx context.Context, id DatabaseAccountId, input RegionForOnlineOffline) error {
	result, err := c.DatabaseAccountsOfflineRegion(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabaseAccountsOfflineRegion: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabaseAccountsOfflineRegion: %+v", err)
	}

	return nil
}

// preparerForDatabaseAccountsOfflineRegion prepares the DatabaseAccountsOfflineRegion request.
func (c CosmosDBClient) preparerForDatabaseAccountsOfflineRegion(ctx context.Context, id DatabaseAccountId, input RegionForOnlineOffline) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/offlineRegion", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDatabaseAccountsOfflineRegion sends the DatabaseAccountsOfflineRegion request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForDatabaseAccountsOfflineRegion(ctx context.Context, req *http.Request) (future DatabaseAccountsOfflineRegionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
