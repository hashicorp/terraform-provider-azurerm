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

type TransferDisksOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TransferDisks ...
func (c VirtualMachinesClient) TransferDisks(ctx context.Context, id VirtualMachineId) (result TransferDisksOperationResponse, err error) {
	req, err := c.preparerForTransferDisks(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "TransferDisks", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTransferDisks(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "TransferDisks", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TransferDisksThenPoll performs TransferDisks then polls until it's completed
func (c VirtualMachinesClient) TransferDisksThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.TransferDisks(ctx, id)
	if err != nil {
		return fmt.Errorf("performing TransferDisks: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TransferDisks: %+v", err)
	}

	return nil
}

// preparerForTransferDisks prepares the TransferDisks request.
func (c VirtualMachinesClient) preparerForTransferDisks(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/transferDisks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTransferDisks sends the TransferDisks request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForTransferDisks(ctx context.Context, req *http.Request) (future TransferDisksOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
