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

type ComputeUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ComputeUpdate ...
func (c MachineLearningComputesClient) ComputeUpdate(ctx context.Context, id ComputeId, input ClusterUpdateParameters) (result ComputeUpdateOperationResponse, err error) {
	req, err := c.preparerForComputeUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForComputeUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ComputeUpdateThenPoll performs ComputeUpdate then polls until it's completed
func (c MachineLearningComputesClient) ComputeUpdateThenPoll(ctx context.Context, id ComputeId, input ClusterUpdateParameters) error {
	result, err := c.ComputeUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ComputeUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ComputeUpdate: %+v", err)
	}

	return nil
}

// preparerForComputeUpdate prepares the ComputeUpdate request.
func (c MachineLearningComputesClient) preparerForComputeUpdate(ctx context.Context, id ComputeId, input ClusterUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForComputeUpdate sends the ComputeUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c MachineLearningComputesClient) senderForComputeUpdate(ctx context.Context, req *http.Request) (future ComputeUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
