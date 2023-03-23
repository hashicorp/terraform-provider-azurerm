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

type CassandraClustersCreateUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraClustersCreateUpdate ...
func (c ManagedCassandrasClient) CassandraClustersCreateUpdate(ctx context.Context, id CassandraClusterId, input ClusterResource) (result CassandraClustersCreateUpdateOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersCreateUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersCreateUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraClustersCreateUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersCreateUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraClustersCreateUpdateThenPoll performs CassandraClustersCreateUpdate then polls until it's completed
func (c ManagedCassandrasClient) CassandraClustersCreateUpdateThenPoll(ctx context.Context, id CassandraClusterId, input ClusterResource) error {
	result, err := c.CassandraClustersCreateUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersCreateUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraClustersCreateUpdate: %+v", err)
	}

	return nil
}

// preparerForCassandraClustersCreateUpdate prepares the CassandraClustersCreateUpdate request.
func (c ManagedCassandrasClient) preparerForCassandraClustersCreateUpdate(ctx context.Context, id CassandraClusterId, input ClusterResource) (*http.Request, error) {
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

// senderForCassandraClustersCreateUpdate sends the CassandraClustersCreateUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraClustersCreateUpdate(ctx context.Context, req *http.Request) (future CassandraClustersCreateUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
