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

type RunCommandOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RunCommand ...
func (c VirtualMachinesClient) RunCommand(ctx context.Context, id VirtualMachineId, input RunCommandInput) (result RunCommandOperationResponse, err error) {
	req, err := c.preparerForRunCommand(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "RunCommand", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRunCommand(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "RunCommand", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RunCommandThenPoll performs RunCommand then polls until it's completed
func (c VirtualMachinesClient) RunCommandThenPoll(ctx context.Context, id VirtualMachineId, input RunCommandInput) error {
	result, err := c.RunCommand(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing RunCommand: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RunCommand: %+v", err)
	}

	return nil
}

// preparerForRunCommand prepares the RunCommand request.
func (c VirtualMachinesClient) preparerForRunCommand(ctx context.Context, id VirtualMachineId, input RunCommandInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/runCommand", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRunCommand sends the RunCommand request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForRunCommand(ctx context.Context, req *http.Request) (future RunCommandOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
