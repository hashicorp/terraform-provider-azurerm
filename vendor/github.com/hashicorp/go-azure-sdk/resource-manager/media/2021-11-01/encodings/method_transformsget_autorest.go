package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Transform
}

// TransformsGet ...
func (c EncodingsClient) TransformsGet(ctx context.Context, id TransformId) (result TransformsGetOperationResponse, err error) {
	req, err := c.preparerForTransformsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTransformsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTransformsGet prepares the TransformsGet request.
func (c EncodingsClient) preparerForTransformsGet(ctx context.Context, id TransformId) (*http.Request, error) {
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

// responderForTransformsGet handles the response to the TransformsGet request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForTransformsGet(resp *http.Response) (result TransformsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
