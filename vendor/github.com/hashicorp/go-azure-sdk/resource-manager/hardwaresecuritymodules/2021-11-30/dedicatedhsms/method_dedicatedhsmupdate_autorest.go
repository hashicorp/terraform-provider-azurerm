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

type DedicatedHsmUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DedicatedHsmUpdate ...
func (c DedicatedHsmsClient) DedicatedHsmUpdate(ctx context.Context, id DedicatedHSMId, input DedicatedHsmPatchParameters) (result DedicatedHsmUpdateOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDedicatedHsmUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DedicatedHsmUpdateThenPoll performs DedicatedHsmUpdate then polls until it's completed
func (c DedicatedHsmsClient) DedicatedHsmUpdateThenPoll(ctx context.Context, id DedicatedHSMId, input DedicatedHsmPatchParameters) error {
	result, err := c.DedicatedHsmUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DedicatedHsmUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DedicatedHsmUpdate: %+v", err)
	}

	return nil
}

// preparerForDedicatedHsmUpdate prepares the DedicatedHsmUpdate request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmUpdate(ctx context.Context, id DedicatedHSMId, input DedicatedHsmPatchParameters) (*http.Request, error) {
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

// senderForDedicatedHsmUpdate sends the DedicatedHsmUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c DedicatedHsmsClient) senderForDedicatedHsmUpdate(ctx context.Context, req *http.Request) (future DedicatedHsmUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
