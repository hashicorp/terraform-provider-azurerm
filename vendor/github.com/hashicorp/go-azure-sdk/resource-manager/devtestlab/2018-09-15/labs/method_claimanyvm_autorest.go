package labs

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

type ClaimAnyVMOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ClaimAnyVM ...
func (c LabsClient) ClaimAnyVM(ctx context.Context, id LabId) (result ClaimAnyVMOperationResponse, err error) {
	req, err := c.preparerForClaimAnyVM(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ClaimAnyVM", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForClaimAnyVM(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ClaimAnyVM", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ClaimAnyVMThenPoll performs ClaimAnyVM then polls until it's completed
func (c LabsClient) ClaimAnyVMThenPoll(ctx context.Context, id LabId) error {
	result, err := c.ClaimAnyVM(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ClaimAnyVM: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ClaimAnyVM: %+v", err)
	}

	return nil
}

// preparerForClaimAnyVM prepares the ClaimAnyVM request.
func (c LabsClient) preparerForClaimAnyVM(ctx context.Context, id LabId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/claimAnyVm", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForClaimAnyVM sends the ClaimAnyVM request. The method will close the
// http.Response Body if it receives an error.
func (c LabsClient) senderForClaimAnyVM(ctx context.Context, req *http.Request) (future ClaimAnyVMOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
