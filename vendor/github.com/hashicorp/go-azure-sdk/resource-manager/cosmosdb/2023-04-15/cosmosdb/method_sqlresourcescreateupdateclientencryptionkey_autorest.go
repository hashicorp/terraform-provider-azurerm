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

type SqlResourcesCreateUpdateClientEncryptionKeyOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SqlResourcesCreateUpdateClientEncryptionKey ...
func (c CosmosDBClient) SqlResourcesCreateUpdateClientEncryptionKey(ctx context.Context, id ClientEncryptionKeyId, input ClientEncryptionKeyCreateUpdateParameters) (result SqlResourcesCreateUpdateClientEncryptionKeyOperationResponse, err error) {
	req, err := c.preparerForSqlResourcesCreateUpdateClientEncryptionKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateClientEncryptionKey", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSqlResourcesCreateUpdateClientEncryptionKey(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "SqlResourcesCreateUpdateClientEncryptionKey", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SqlResourcesCreateUpdateClientEncryptionKeyThenPoll performs SqlResourcesCreateUpdateClientEncryptionKey then polls until it's completed
func (c CosmosDBClient) SqlResourcesCreateUpdateClientEncryptionKeyThenPoll(ctx context.Context, id ClientEncryptionKeyId, input ClientEncryptionKeyCreateUpdateParameters) error {
	result, err := c.SqlResourcesCreateUpdateClientEncryptionKey(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SqlResourcesCreateUpdateClientEncryptionKey: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SqlResourcesCreateUpdateClientEncryptionKey: %+v", err)
	}

	return nil
}

// preparerForSqlResourcesCreateUpdateClientEncryptionKey prepares the SqlResourcesCreateUpdateClientEncryptionKey request.
func (c CosmosDBClient) preparerForSqlResourcesCreateUpdateClientEncryptionKey(ctx context.Context, id ClientEncryptionKeyId, input ClientEncryptionKeyCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSqlResourcesCreateUpdateClientEncryptionKey sends the SqlResourcesCreateUpdateClientEncryptionKey request. The method will close the
// http.Response Body if it receives an error.
func (c CosmosDBClient) senderForSqlResourcesCreateUpdateClientEncryptionKey(ctx context.Context, req *http.Request) (future SqlResourcesCreateUpdateClientEncryptionKeyOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
