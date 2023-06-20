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

type RotateNsxtPasswordOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RotateNsxtPassword ...
func (c PrivateCloudsClient) RotateNsxtPassword(ctx context.Context, id PrivateCloudId) (result RotateNsxtPasswordOperationResponse, err error) {
	req, err := c.preparerForRotateNsxtPassword(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "RotateNsxtPassword", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRotateNsxtPassword(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privateclouds.PrivateCloudsClient", "RotateNsxtPassword", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RotateNsxtPasswordThenPoll performs RotateNsxtPassword then polls until it's completed
func (c PrivateCloudsClient) RotateNsxtPasswordThenPoll(ctx context.Context, id PrivateCloudId) error {
	result, err := c.RotateNsxtPassword(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RotateNsxtPassword: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RotateNsxtPassword: %+v", err)
	}

	return nil
}

// preparerForRotateNsxtPassword prepares the RotateNsxtPassword request.
func (c PrivateCloudsClient) preparerForRotateNsxtPassword(ctx context.Context, id PrivateCloudId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rotateNsxtPassword", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRotateNsxtPassword sends the RotateNsxtPassword request. The method will close the
// http.Response Body if it receives an error.
func (c PrivateCloudsClient) senderForRotateNsxtPassword(ctx context.Context, req *http.Request) (future RotateNsxtPasswordOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
