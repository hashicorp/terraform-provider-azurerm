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

type InstancesCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// InstancesCreate ...
func (c DeviceupdatesClient) InstancesCreate(ctx context.Context, id InstanceId, input Instance) (result InstancesCreateOperationResponse, err error) {
	req, err := c.preparerForInstancesCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForInstancesCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// InstancesCreateThenPoll performs InstancesCreate then polls until it's completed
func (c DeviceupdatesClient) InstancesCreateThenPoll(ctx context.Context, id InstanceId, input Instance) error {
	result, err := c.InstancesCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing InstancesCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after InstancesCreate: %+v", err)
	}

	return nil
}

// preparerForInstancesCreate prepares the InstancesCreate request.
func (c DeviceupdatesClient) preparerForInstancesCreate(ctx context.Context, id InstanceId, input Instance) (*http.Request, error) {
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

// senderForInstancesCreate sends the InstancesCreate request. The method will close the
// http.Response Body if it receives an error.
func (c DeviceupdatesClient) senderForInstancesCreate(ctx context.Context, req *http.Request) (future InstancesCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
