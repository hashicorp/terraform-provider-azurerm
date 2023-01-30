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

type AddDataDiskOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AddDataDisk ...
func (c VirtualMachinesClient) AddDataDisk(ctx context.Context, id VirtualMachineId, input DataDiskProperties) (result AddDataDiskOperationResponse, err error) {
	req, err := c.preparerForAddDataDisk(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "AddDataDisk", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAddDataDisk(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "AddDataDisk", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AddDataDiskThenPoll performs AddDataDisk then polls until it's completed
func (c VirtualMachinesClient) AddDataDiskThenPoll(ctx context.Context, id VirtualMachineId, input DataDiskProperties) error {
	result, err := c.AddDataDisk(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AddDataDisk: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AddDataDisk: %+v", err)
	}

	return nil
}

// preparerForAddDataDisk prepares the AddDataDisk request.
func (c VirtualMachinesClient) preparerForAddDataDisk(ctx context.Context, id VirtualMachineId, input DataDiskProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/addDataDisk", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAddDataDisk sends the AddDataDisk request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForAddDataDisk(ctx context.Context, req *http.Request) (future AddDataDiskOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
