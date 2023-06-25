package privateclouds

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

type RotateVcenterPasswordOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RotateVcenterPassword ...
func (c PrivateCloudsClient) RotateVcenterPassword(ctx context.Context, id PrivateCloudId) (result RotateVcenterPasswordOperationResponse, err error) {
	req, err := c.preparerForRotateVcenterPassword(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "RotateVcenterPassword", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRotateVcenterPassword(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "RotateVcenterPassword", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RotateVcenterPasswordThenPoll performs RotateVcenterPassword then polls until it's completed
func (c PrivateCloudsClient) RotateVcenterPasswordThenPoll(ctx context.Context, id PrivateCloudId) error {
	result, err := c.RotateVcenterPassword(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RotateVcenterPassword: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RotateVcenterPassword: %+v", err)
	}

	return nil
}

// preparerForRotateVcenterPassword prepares the RotateVcenterPassword request.
func (c PrivateCloudsClient) preparerForRotateVcenterPassword(ctx context.Context, id PrivateCloudId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rotateVcenterPassword", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRotateVcenterPassword sends the RotateVcenterPassword request. The method will close the
// http.Response Body if it receives an error.
func (c PrivateCloudsClient) senderForRotateVcenterPassword(ctx context.Context, req *http.Request) (future RotateVcenterPasswordOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
