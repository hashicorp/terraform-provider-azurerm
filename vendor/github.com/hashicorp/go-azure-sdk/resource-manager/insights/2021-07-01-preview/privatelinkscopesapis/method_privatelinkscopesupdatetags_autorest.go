package privatelinkscopesapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesUpdateTagsOperationResponse struct {
	HttpResponse *http.Response
	Model        *AzureMonitorPrivateLinkScope
}

// PrivateLinkScopesUpdateTags ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesUpdateTags(ctx context.Context, id PrivateLinkScopeId, input TagsResource) (result PrivateLinkScopesUpdateTagsOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesUpdateTags(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesUpdateTags", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesUpdateTags", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkScopesUpdateTags(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesUpdateTags", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkScopesUpdateTags prepares the PrivateLinkScopesUpdateTags request.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesUpdateTags(ctx context.Context, id PrivateLinkScopeId, input TagsResource) (*http.Request, error) {
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

// responderForPrivateLinkScopesUpdateTags handles the response to the PrivateLinkScopesUpdateTags request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesAPIsClient) responderForPrivateLinkScopesUpdateTags(resp *http.Response) (result PrivateLinkScopesUpdateTagsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
