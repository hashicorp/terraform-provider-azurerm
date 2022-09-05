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

type CassandraClustersUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraClustersUpdate ...
func (c ManagedCassandrasClient) CassandraClustersUpdate(ctx context.Context, id CassandraClusterId, input ClusterResource) (result CassandraClustersUpdateOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraClustersUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraClustersUpdateThenPoll performs CassandraClustersUpdate then polls until it's completed
func (c ManagedCassandrasClient) CassandraClustersUpdateThenPoll(ctx context.Context, id CassandraClusterId, input ClusterResource) error {
	result, err := c.CassandraClustersUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraClustersUpdate: %+v", err)
	}

	return nil
}

// preparerForCassandraClustersUpdate prepares the CassandraClustersUpdate request.
func (c ManagedCassandrasClient) preparerForCassandraClustersUpdate(ctx context.Context, id CassandraClusterId, input ClusterResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraClustersUpdate sends the CassandraClustersUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraClustersUpdate(ctx context.Context, req *http.Request) (future CassandraClustersUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
