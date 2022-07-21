package dedicatedhsms

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

type DedicatedHsmDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DedicatedHsmDelete ...
func (c DedicatedHsmsClient) DedicatedHsmDelete(ctx context.Context, id DedicatedHSMId) (result DedicatedHsmDeleteOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDedicatedHsmDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DedicatedHsmDeleteThenPoll performs DedicatedHsmDelete then polls until it's completed
func (c DedicatedHsmsClient) DedicatedHsmDeleteThenPoll(ctx context.Context, id DedicatedHSMId) error {
	result, err := c.DedicatedHsmDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DedicatedHsmDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DedicatedHsmDelete: %+v", err)
	}

	return nil
}

// preparerForDedicatedHsmDelete prepares the DedicatedHsmDelete request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmDelete(ctx context.Context, id DedicatedHSMId) (*http.Request, error) {
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

// senderForDedicatedHsmDelete sends the DedicatedHsmDelete request. The method will close the
// http.Response Body if it receives an error.
func (c DedicatedHsmsClient) senderForDedicatedHsmDelete(ctx context.Context, req *http.Request) (future DedicatedHsmDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
