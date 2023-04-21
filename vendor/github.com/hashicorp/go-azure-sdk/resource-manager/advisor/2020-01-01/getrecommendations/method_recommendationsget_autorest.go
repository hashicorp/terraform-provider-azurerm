package getrecommendations

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecommendationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ResourceRecommendationBase
}

// RecommendationsGet ...
func (c GetRecommendationsClient) RecommendationsGet(ctx context.Context, id ScopedRecommendationId) (result RecommendationsGetOperationResponse, err error) {
	req, err := c.preparerForRecommendationsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRecommendationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRecommendationsGet prepares the RecommendationsGet request.
func (c GetRecommendationsClient) preparerForRecommendationsGet(ctx context.Context, id ScopedRecommendationId) (*http.Request, error) {
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

// responderForRecommendationsGet handles the response to the RecommendationsGet request. The method always
// closes the http.Response Body.
func (c GetRecommendationsClient) responderForRecommendationsGet(resp *http.Response) (result RecommendationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
