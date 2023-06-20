package machinelearningcomputes

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

type ComputeStartOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ComputeStart ...
func (c MachineLearningComputesClient) ComputeStart(ctx context.Context, id ComputeId) (result ComputeStartOperationResponse, err error) {
	req, err := c.preparerForComputeStart(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeStart", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForComputeStart(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeStart", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ComputeStartThenPoll performs ComputeStart then polls until it's completed
func (c MachineLearningComputesClient) ComputeStartThenPoll(ctx context.Context, id ComputeId) error {
	result, err := c.ComputeStart(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ComputeStart: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ComputeStart: %+v", err)
	}

	return nil
}

// preparerForComputeStart prepares the ComputeStart request.
func (c MachineLearningComputesClient) preparerForComputeStart(ctx context.Context, id ComputeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/start", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForComputeStart sends the ComputeStart request. The method will close the
// http.Response Body if it receives an error.
func (c MachineLearningComputesClient) senderForComputeStart(ctx context.Context, req *http.Request) (future ComputeStartOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
