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

type InstallPatchesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// InstallPatches ...
func (c VirtualMachinesClient) InstallPatches(ctx context.Context, id VirtualMachineId, input VirtualMachineInstallPatchesParameters) (result InstallPatchesOperationResponse, err error) {
	req, err := c.preparerForInstallPatches(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "InstallPatches", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForInstallPatches(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "InstallPatches", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// InstallPatchesThenPoll performs InstallPatches then polls until it's completed
func (c VirtualMachinesClient) InstallPatchesThenPoll(ctx context.Context, id VirtualMachineId, input VirtualMachineInstallPatchesParameters) error {
	result, err := c.InstallPatches(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing InstallPatches: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after InstallPatches: %+v", err)
	}

	return nil
}

// preparerForInstallPatches prepares the InstallPatches request.
func (c VirtualMachinesClient) preparerForInstallPatches(ctx context.Context, id VirtualMachineId, input VirtualMachineInstallPatchesParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/installPatches", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForInstallPatches sends the InstallPatches request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForInstallPatches(ctx context.Context, req *http.Request) (future InstallPatchesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
