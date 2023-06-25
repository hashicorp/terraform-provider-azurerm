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

type DetachDataDiskOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DetachDataDisk ...
func (c VirtualMachinesClient) DetachDataDisk(ctx context.Context, id VirtualMachineId, input DetachDataDiskProperties) (result DetachDataDiskOperationResponse, err error) {
	req, err := c.preparerForDetachDataDisk(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "DetachDataDisk", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDetachDataDisk(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "DetachDataDisk", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DetachDataDiskThenPoll performs DetachDataDisk then polls until it's completed
func (c VirtualMachinesClient) DetachDataDiskThenPoll(ctx context.Context, id VirtualMachineId, input DetachDataDiskProperties) error {
	result, err := c.DetachDataDisk(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DetachDataDisk: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DetachDataDisk: %+v", err)
	}

	return nil
}

// preparerForDetachDataDisk prepares the DetachDataDisk request.
func (c VirtualMachinesClient) preparerForDetachDataDisk(ctx context.Context, id VirtualMachineId, input DetachDataDiskProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/detachDataDisk", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDetachDataDisk sends the DetachDataDisk request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForDetachDataDisk(ctx context.Context, req *http.Request) (future DetachDataDiskOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
