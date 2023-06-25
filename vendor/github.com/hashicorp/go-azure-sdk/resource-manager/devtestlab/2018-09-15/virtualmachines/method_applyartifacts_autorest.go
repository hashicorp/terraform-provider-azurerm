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

type ApplyArtifactsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ApplyArtifacts ...
func (c VirtualMachinesClient) ApplyArtifacts(ctx context.Context, id VirtualMachineId, input ApplyArtifactsRequest) (result ApplyArtifactsOperationResponse, err error) {
	req, err := c.preparerForApplyArtifacts(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ApplyArtifacts", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForApplyArtifacts(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ApplyArtifacts", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ApplyArtifactsThenPoll performs ApplyArtifacts then polls until it's completed
func (c VirtualMachinesClient) ApplyArtifactsThenPoll(ctx context.Context, id VirtualMachineId, input ApplyArtifactsRequest) error {
	result, err := c.ApplyArtifacts(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ApplyArtifacts: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ApplyArtifacts: %+v", err)
	}

	return nil
}

// preparerForApplyArtifacts prepares the ApplyArtifacts request.
func (c VirtualMachinesClient) preparerForApplyArtifacts(ctx context.Context, id VirtualMachineId, input ApplyArtifactsRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/applyArtifacts", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForApplyArtifacts sends the ApplyArtifacts request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForApplyArtifacts(ctx context.Context, req *http.Request) (future ApplyArtifactsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
