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

type CassandraClustersStartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraClustersStart ...
func (c ManagedCassandrasClient) CassandraClustersStart(ctx context.Context, id CassandraClusterId) (result CassandraClustersStartOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersStart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersStart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraClustersStart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersStart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraClustersStartThenPoll performs CassandraClustersStart then polls until it's completed
func (c ManagedCassandrasClient) CassandraClustersStartThenPoll(ctx context.Context, id CassandraClusterId) error {
	result, err := c.CassandraClustersStart(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersStart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraClustersStart: %+v", err)
	}

	return nil
}

// preparerForCassandraClustersStart prepares the CassandraClustersStart request.
func (c ManagedCassandrasClient) preparerForCassandraClustersStart(ctx context.Context, id CassandraClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/start", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraClustersStart sends the CassandraClustersStart request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraClustersStart(ctx context.Context, req *http.Request) (future CassandraClustersStartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
