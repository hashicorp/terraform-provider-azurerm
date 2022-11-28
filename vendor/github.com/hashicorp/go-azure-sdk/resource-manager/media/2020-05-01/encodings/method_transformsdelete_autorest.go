package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// TransformsDelete ...
func (c EncodingsClient) TransformsDelete(ctx context.Context, id TransformId) (result TransformsDeleteOperationResponse, err error) {
	req, err := c.preparerForTransformsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTransformsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTransformsDelete prepares the TransformsDelete request.
func (c EncodingsClient) preparerForTransformsDelete(ctx context.Context, id TransformId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTransformsDelete handles the response to the TransformsDelete request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForTransformsDelete(resp *http.Response) (result TransformsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
