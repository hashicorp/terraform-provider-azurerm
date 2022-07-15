package machinelearningcomputes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *ComputeSecrets
}

// ComputeListKeys ...
func (c MachineLearningComputesClient) ComputeListKeys(ctx context.Context, id ComputeId) (result ComputeListKeysOperationResponse, err error) {
	req, err := c.preparerForComputeListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForComputeListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForComputeListKeys prepares the ComputeListKeys request.
func (c MachineLearningComputesClient) preparerForComputeListKeys(ctx context.Context, id ComputeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForComputeListKeys handles the response to the ComputeListKeys request. The method always
// closes the http.Response Body.
func (c MachineLearningComputesClient) responderForComputeListKeys(resp *http.Response) (result ComputeListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
