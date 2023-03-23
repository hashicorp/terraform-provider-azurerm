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

type CassandraClustersInvokeCommandOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraClustersInvokeCommand ...
func (c ManagedCassandrasClient) CassandraClustersInvokeCommand(ctx context.Context, id CassandraClusterId, input CommandPostBody) (result CassandraClustersInvokeCommandOperationResponse, err error) {
	req, err := c.preparerForCassandraClustersInvokeCommand(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersInvokeCommand", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraClustersInvokeCommand(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraClustersInvokeCommand", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraClustersInvokeCommandThenPoll performs CassandraClustersInvokeCommand then polls until it's completed
func (c ManagedCassandrasClient) CassandraClustersInvokeCommandThenPoll(ctx context.Context, id CassandraClusterId, input CommandPostBody) error {
	result, err := c.CassandraClustersInvokeCommand(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraClustersInvokeCommand: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraClustersInvokeCommand: %+v", err)
	}

	return nil
}

// preparerForCassandraClustersInvokeCommand prepares the CassandraClustersInvokeCommand request.
func (c ManagedCassandrasClient) preparerForCassandraClustersInvokeCommand(ctx context.Context, id CassandraClusterId, input CommandPostBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/invokeCommand", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCassandraClustersInvokeCommand sends the CassandraClustersInvokeCommand request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraClustersInvokeCommand(ctx context.Context, req *http.Request) (future CassandraClustersInvokeCommandOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
