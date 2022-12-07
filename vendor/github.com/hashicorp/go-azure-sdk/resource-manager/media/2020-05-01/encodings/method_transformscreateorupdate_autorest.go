package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Transform
}

// TransformsCreateOrUpdate ...
func (c EncodingsClient) TransformsCreateOrUpdate(ctx context.Context, id TransformId, input Transform) (result TransformsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForTransformsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTransformsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTransformsCreateOrUpdate prepares the TransformsCreateOrUpdate request.
func (c EncodingsClient) preparerForTransformsCreateOrUpdate(ctx context.Context, id TransformId, input Transform) (*http.Request, error) {
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

// responderForTransformsCreateOrUpdate handles the response to the TransformsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForTransformsCreateOrUpdate(resp *http.Response) (result TransformsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
