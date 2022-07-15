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

type ComputeStopOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ComputeStop ...
func (c MachineLearningComputesClient) ComputeStop(ctx context.Context, id ComputeId) (result ComputeStopOperationResponse, err error) {
	req, err := c.preparerForComputeStop(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeStop", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForComputeStop(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeStop", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ComputeStopThenPoll performs ComputeStop then polls until it's completed
func (c MachineLearningComputesClient) ComputeStopThenPoll(ctx context.Context, id ComputeId) error {
	result, err := c.ComputeStop(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ComputeStop: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ComputeStop: %+v", err)
	}

	return nil
}

// preparerForComputeStop prepares the ComputeStop request.
func (c MachineLearningComputesClient) preparerForComputeStop(ctx context.Context, id ComputeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stop", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForComputeStop sends the ComputeStop request. The method will close the
// http.Response Body if it receives an error.
func (c MachineLearningComputesClient) senderForComputeStop(ctx context.Context, req *http.Request) (future ComputeStopOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
