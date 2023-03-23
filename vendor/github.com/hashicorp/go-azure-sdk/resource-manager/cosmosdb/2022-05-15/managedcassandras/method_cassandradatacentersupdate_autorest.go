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

type CassandraDataCentersUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraDataCentersUpdate ...
func (c ManagedCassandrasClient) CassandraDataCentersUpdate(ctx context.Context, id DataCenterId, input DataCenterResource) (result CassandraDataCentersUpdateOperationResponse, err error) {
	req, err := c.preparerForCassandraDataCentersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraDataCentersUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraDataCentersUpdateThenPoll performs CassandraDataCentersUpdate then polls until it's completed
func (c ManagedCassandrasClient) CassandraDataCentersUpdateThenPoll(ctx context.Context, id DataCenterId, input DataCenterResource) error {
	result, err := c.CassandraDataCentersUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraDataCentersUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraDataCentersUpdate: %+v", err)
	}

	return nil
}

// preparerForCassandraDataCentersUpdate prepares the CassandraDataCentersUpdate request.
func (c ManagedCassandrasClient) preparerForCassandraDataCentersUpdate(ctx context.Context, id DataCenterId, input DataCenterResource) (*http.Request, error) {
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

// senderForCassandraDataCentersUpdate sends the CassandraDataCentersUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraDataCentersUpdate(ctx context.Context, req *http.Request) (future CassandraDataCentersUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
