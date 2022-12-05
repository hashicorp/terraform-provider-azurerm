package encodings

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Transform
}

// TransformsUpdate ...
func (c EncodingsClient) TransformsUpdate(ctx context.Context, id TransformId, input Transform) (result TransformsUpdateOperationResponse, err error) {
	req, err := c.preparerForTransformsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTransformsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTransformsUpdate prepares the TransformsUpdate request.
func (c EncodingsClient) preparerForTransformsUpdate(ctx context.Context, id TransformId, input Transform) (*http.Request, error) {
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

// responderForTransformsUpdate handles the response to the TransformsUpdate request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForTransformsUpdate(resp *http.Response) (result TransformsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
