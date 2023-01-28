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

type ResizeOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Resize ...
func (c VirtualMachinesClient) Resize(ctx context.Context, id VirtualMachineId, input ResizeLabVirtualMachineProperties) (result ResizeOperationResponse, err error) {
	req, err := c.preparerForResize(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Resize", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResize(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Resize", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResizeThenPoll performs Resize then polls until it's completed
func (c VirtualMachinesClient) ResizeThenPoll(ctx context.Context, id VirtualMachineId, input ResizeLabVirtualMachineProperties) error {
	result, err := c.Resize(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Resize: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Resize: %+v", err)
	}

	return nil
}

// preparerForResize prepares the Resize request.
func (c VirtualMachinesClient) preparerForResize(ctx context.Context, id VirtualMachineId, input ResizeLabVirtualMachineProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resize", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResize sends the Resize request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForResize(ctx context.Context, req *http.Request) (future ResizeOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
