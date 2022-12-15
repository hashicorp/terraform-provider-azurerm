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

type ReapplyOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Reapply ...
func (c VirtualMachinesClient) Reapply(ctx context.Context, id VirtualMachineId) (result ReapplyOperationResponse, err error) {
	req, err := c.preparerForReapply(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Reapply", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReapply(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "Reapply", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReapplyThenPoll performs Reapply then polls until it's completed
func (c VirtualMachinesClient) ReapplyThenPoll(ctx context.Context, id VirtualMachineId) error {
	result, err := c.Reapply(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Reapply: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Reapply: %+v", err)
	}

	return nil
}

// preparerForReapply prepares the Reapply request.
func (c VirtualMachinesClient) preparerForReapply(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reapply", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReapply sends the Reapply request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForReapply(ctx context.Context, req *http.Request) (future ReapplyOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
