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

type CassandraDataCentersCreateUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CassandraDataCentersCreateUpdate ...
func (c ManagedCassandrasClient) CassandraDataCentersCreateUpdate(ctx context.Context, id DataCenterId, input DataCenterResource) (result CassandraDataCentersCreateUpdateOperationResponse, err error) {
	req, err := c.preparerForCassandraDataCentersCreateUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersCreateUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCassandraDataCentersCreateUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedcassandras.ManagedCassandrasClient", "CassandraDataCentersCreateUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CassandraDataCentersCreateUpdateThenPoll performs CassandraDataCentersCreateUpdate then polls until it's completed
func (c ManagedCassandrasClient) CassandraDataCentersCreateUpdateThenPoll(ctx context.Context, id DataCenterId, input DataCenterResource) error {
	result, err := c.CassandraDataCentersCreateUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CassandraDataCentersCreateUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CassandraDataCentersCreateUpdate: %+v", err)
	}

	return nil
}

// preparerForCassandraDataCentersCreateUpdate prepares the CassandraDataCentersCreateUpdate request.
func (c ManagedCassandrasClient) preparerForCassandraDataCentersCreateUpdate(ctx context.Context, id DataCenterId, input DataCenterResource) (*http.Request, error) {
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

// senderForCassandraDataCentersCreateUpdate sends the CassandraDataCentersCreateUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedCassandrasClient) senderForCassandraDataCentersCreateUpdate(ctx context.Context, req *http.Request) (future CassandraDataCentersCreateUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
