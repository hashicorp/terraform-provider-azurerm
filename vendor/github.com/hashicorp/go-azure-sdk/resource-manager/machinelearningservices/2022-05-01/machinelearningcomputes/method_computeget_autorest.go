package machinelearningcomputes

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ComputeResource
}

// ComputeGet ...
func (c MachineLearningComputesClient) ComputeGet(ctx context.Context, id ComputeId) (result ComputeGetOperationResponse, err error) {
	req, err := c.preparerForComputeGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForComputeGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForComputeGet prepares the ComputeGet request.
func (c MachineLearningComputesClient) preparerForComputeGet(ctx context.Context, id ComputeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForComputeGet handles the response to the ComputeGet request. The method always
// closes the http.Response Body.
func (c MachineLearningComputesClient) responderForComputeGet(resp *http.Response) (result ComputeGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
