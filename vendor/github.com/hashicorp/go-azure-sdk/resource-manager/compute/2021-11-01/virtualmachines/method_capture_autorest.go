package virtualmachines

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

type CaptureOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Capture ...
func (c VirtualMachinesClient) Capture(ctx context.Context, id VirtualMachineId, input VirtualMachineCaptureParameters) (result CaptureOperationResponse, err error) {
	req, err := c.preparerForCapture(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Capture", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCapture(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Capture", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CaptureThenPoll performs Capture then polls until it's completed
func (c VirtualMachinesClient) CaptureThenPoll(ctx context.Context, id VirtualMachineId, input VirtualMachineCaptureParameters) error {
	result, err := c.Capture(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Capture: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Capture: %+v", err)
	}

	return nil
}

// preparerForCapture prepares the Capture request.
func (c VirtualMachinesClient) preparerForCapture(ctx context.Context, id VirtualMachineId, input VirtualMachineCaptureParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/capture", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCapture sends the Capture request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForCapture(ctx context.Context, req *http.Request) (future CaptureOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
