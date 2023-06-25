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

type PerformMaintenanceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PerformMaintenance ...
func (c VirtualMachinesClient) PerformMaintenance(ctx context.Context, id VirtualMachineId) (result PerformMaintenanceOperationResponse, err error) {
	req, err := c.preparerForPerformMaintenance(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "PerformMaintenance", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPerformMaintenance(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "PerformMaintenance", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PerformMaintenanceThenPoll performs PerformMaintenance then polls until it's completed
func (c VirtualMachinesClient) PerformMaintenanceThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.PerformMaintenance(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PerformMaintenance: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PerformMaintenance: %+v", err)
	}

	return nil
}

// preparerForPerformMaintenance prepares the PerformMaintenance request.
func (c VirtualMachinesClient) preparerForPerformMaintenance(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/performMaintenance", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPerformMaintenance sends the PerformMaintenance request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForPerformMaintenance(ctx context.Context, req *http.Request) (future PerformMaintenanceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
