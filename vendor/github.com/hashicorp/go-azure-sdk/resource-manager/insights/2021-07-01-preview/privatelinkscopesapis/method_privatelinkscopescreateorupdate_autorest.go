package privatelinkscopesapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AzureMonitorPrivateLinkScope
}

// PrivateLinkScopesCreateOrUpdate ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesCreateOrUpdate(ctx context.Context, id PrivateLinkScopeId, input AzureMonitorPrivateLinkScope) (result PrivateLinkScopesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkScopesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkScopesCreateOrUpdate prepares the PrivateLinkScopesCreateOrUpdate request.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesCreateOrUpdate(ctx context.Context, id PrivateLinkScopeId, input AzureMonitorPrivateLinkScope) (*http.Request, error) {
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

// responderForPrivateLinkScopesCreateOrUpdate handles the response to the PrivateLinkScopesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesAPIsClient) responderForPrivateLinkScopesCreateOrUpdate(resp *http.Response) (result PrivateLinkScopesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
