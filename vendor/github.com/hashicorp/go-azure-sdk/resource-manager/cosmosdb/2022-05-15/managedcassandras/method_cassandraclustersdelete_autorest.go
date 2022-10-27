package managedcassandras

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

type CassandraClustersDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraClustersDelete ...
func (c ManagedCassandrasClient) CassandraClustersDelete(ctx context.Context, id CassandraClusterId) (result CassandraClustersDeleteOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraClustersDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraClustersDeleteThenPoll performs CassandraClustersDelete then polls until it's completed
func (c ManagedCassandrasClient) CassandraClustersDeleteThenPoll(ctx context.Context, id CassandraClusterId) error {
	result, err := c.CassandraClustersDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraClustersDelete: %+v", err)
	}

	return nil
}

// preparerForCassandraClustersDelete prepares the CassandraClustersDelete request.
func (c ManagedCassandrasClient) preparerForCassandraClustersDelete(ctx context.Context, id CassandraClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraClustersDelete sends the CassandraClustersDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraClustersDelete(ctx context.Context, req *http.Request) (future CassandraClustersDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
