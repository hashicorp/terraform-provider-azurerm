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

type DedicatedHsmCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DedicatedHsmCreateOrUpdate ...
func (c DedicatedHsmsClient) DedicatedHsmCreateOrUpdate(ctx context.Context, id DedicatedHSMId, input DedicatedHsm) (result DedicatedHsmCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDedicatedHsmCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DedicatedHsmCreateOrUpdateThenPoll performs DedicatedHsmCreateOrUpdate then polls until it's completed
func (c DedicatedHsmsClient) DedicatedHsmCreateOrUpdateThenPoll(ctx context.Context, id DedicatedHSMId, input DedicatedHsm) error {
	result, err := c.DedicatedHsmCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DedicatedHsmCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DedicatedHsmCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForDedicatedHsmCreateOrUpdate prepares the DedicatedHsmCreateOrUpdate request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmCreateOrUpdate(ctx context.Context, id DedicatedHSMId, input DedicatedHsm) (*http.Request, error) {
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

// senderForDedicatedHsmCreateOrUpdate sends the DedicatedHsmCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c DedicatedHsmsClient) senderForDedicatedHsmCreateOrUpdate(ctx context.Context, req *http.Request) (future DedicatedHsmCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
