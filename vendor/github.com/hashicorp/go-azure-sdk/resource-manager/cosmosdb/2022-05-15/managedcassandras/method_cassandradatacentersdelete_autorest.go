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

type CassandraDataCentersDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraDataCentersDelete ...
func (c ManagedCassandrasClient) CassandraDataCentersDelete(ctx context.Context, id DataCenterId) (result CassandraDataCentersDeleteOperationResponse, err error) {
	req, err := c.preparerForCassandraDataCentersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraDataCentersDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraDataCentersDeleteThenPoll performs CassandraDataCentersDelete then polls until it's completed
func (c ManagedCassandrasClient) CassandraDataCentersDeleteThenPoll(ctx context.Context, id DataCenterId) error {
	result, err := c.CassandraDataCentersDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CassandraDataCentersDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraDataCentersDelete: %+v", err)
	}

	return nil
}

// preparerForCassandraDataCentersDelete prepares the CassandraDataCentersDelete request.
func (c ManagedCassandrasClient) preparerForCassandraDataCentersDelete(ctx context.Context, id DataCenterId) (*http.Request, error) {
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

// senderForCassandraDataCentersDelete sends the CassandraDataCentersDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraDataCentersDelete(ctx context.Context, req *http.Request) (future CassandraDataCentersDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
