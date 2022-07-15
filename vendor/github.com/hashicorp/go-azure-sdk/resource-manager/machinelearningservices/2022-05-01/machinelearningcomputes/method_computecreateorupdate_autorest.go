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

type ComputeCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ComputeCreateOrUpdate ...
func (c MachineLearningComputesClient) ComputeCreateOrUpdate(ctx context.Context, id ComputeId, input ComputeResource) (result ComputeCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForComputeCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForComputeCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ComputeCreateOrUpdateThenPoll performs ComputeCreateOrUpdate then polls until it's completed
func (c MachineLearningComputesClient) ComputeCreateOrUpdateThenPoll(ctx context.Context, id ComputeId, input ComputeResource) error {
	result, err := c.ComputeCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ComputeCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ComputeCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForComputeCreateOrUpdate prepares the ComputeCreateOrUpdate request.
func (c MachineLearningComputesClient) preparerForComputeCreateOrUpdate(ctx context.Context, id ComputeId, input ComputeResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForComputeCreateOrUpdate sends the ComputeCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c MachineLearningComputesClient) senderForComputeCreateOrUpdate(ctx context.Context, req *http.Request) (future ComputeCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
