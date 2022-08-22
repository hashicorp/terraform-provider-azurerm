package querypacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksUpdateTagsOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPack
}

// QueryPacksUpdateTags ...
func (c QueryPacksClient) QueryPacksUpdateTags(ctx context.Context, id QueryPackId, input TagsResource) (result QueryPacksUpdateTagsOperationResponse, err error) {
	req, err := c.preparerForQueryPacksUpdateTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksUpdateTags", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksUpdateTags", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueryPacksUpdateTags(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksUpdateTags", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueryPacksUpdateTags prepares the QueryPacksUpdateTags request.
func (c QueryPacksClient) preparerForQueryPacksUpdateTags(ctx context.Context, id QueryPackId, input TagsResource) (*http.Request, error) {
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

// responderForQueryPacksUpdateTags handles the response to the QueryPacksUpdateTags request. The method always
// closes the http.Response Body.
func (c QueryPacksClient) responderForQueryPacksUpdateTags(resp *http.Response) (result QueryPacksUpdateTagsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
