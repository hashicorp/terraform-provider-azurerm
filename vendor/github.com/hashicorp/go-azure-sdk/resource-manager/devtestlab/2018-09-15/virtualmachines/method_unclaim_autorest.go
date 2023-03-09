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

type UnClaimOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UnClaim ...
func (c VirtualMachinesClient) UnClaim(ctx context.Context, id VirtualMachineId) (result UnClaimOperationResponse, err error) {
	req, err := c.preparerForUnClaim(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "UnClaim", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUnClaim(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "UnClaim", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UnClaimThenPoll performs UnClaim then polls until it's completed
func (c VirtualMachinesClient) UnClaimThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.UnClaim(ctx, id)
	if err != nil {
		return fmt.Errorf("performing UnClaim: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UnClaim: %+v", err)
	}

	return nil
}

// preparerForUnClaim prepares the UnClaim request.
func (c VirtualMachinesClient) preparerForUnClaim(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/unClaim", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUnClaim sends the UnClaim request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForUnClaim(ctx context.Context, req *http.Request) (future UnClaimOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
