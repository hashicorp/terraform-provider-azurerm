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

type ConvertToManagedDisksOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConvertToManagedDisks ...
func (c VirtualMachinesClient) ConvertToManagedDisks(ctx context.Context, id VirtualMachineId) (result ConvertToManagedDisksOperationResponse, err error) {
	req, err := c.preparerForConvertToManagedDisks(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ConvertToManagedDisks", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConvertToManagedDisks(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ConvertToManagedDisks", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConvertToManagedDisksThenPoll performs ConvertToManagedDisks then polls until it's completed
func (c VirtualMachinesClient) ConvertToManagedDisksThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.ConvertToManagedDisks(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConvertToManagedDisks: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConvertToManagedDisks: %+v", err)
	}

	return nil
}

// preparerForConvertToManagedDisks prepares the ConvertToManagedDisks request.
func (c VirtualMachinesClient) preparerForConvertToManagedDisks(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/convertToManagedDisks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConvertToManagedDisks sends the ConvertToManagedDisks request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForConvertToManagedDisks(ctx context.Context, req *http.Request) (future ConvertToManagedDisksOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
