package deviceupdates

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

type InstancesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// InstancesDelete ...
func (c DeviceupdatesClient) InstancesDelete(ctx context.Context, id InstanceId) (result InstancesDeleteOperationResponse, err error) {
	req, err := c.preparerForInstancesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForInstancesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// InstancesDeleteThenPoll performs InstancesDelete then polls until it's completed
func (c DeviceupdatesClient) InstancesDeleteThenPoll(ctx context.Context, id InstanceId) error {
	result, err := c.InstancesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing InstancesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after InstancesDelete: %+v", err)
	}

	return nil
}

// preparerForInstancesDelete prepares the InstancesDelete request.
func (c DeviceupdatesClient) preparerForInstancesDelete(ctx context.Context, id InstanceId) (*http.Request, error) {
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

// senderForInstancesDelete sends the InstancesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c DeviceupdatesClient) senderForInstancesDelete(ctx context.Context, req *http.Request) (future InstancesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
