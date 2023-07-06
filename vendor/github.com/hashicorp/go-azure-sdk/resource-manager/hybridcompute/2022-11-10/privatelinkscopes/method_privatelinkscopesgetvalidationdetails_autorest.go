package privatelinkscopes

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesGetValidationDetailsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkScopeValidationDetails
}

// PrivateLinkScopesGetValidationDetails ...
func (c PrivateLinkScopesClient) PrivateLinkScopesGetValidationDetails(ctx context.Context, id PrivateLinkScopeId) (result PrivateLinkScopesGetValidationDetailsOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesGetValidationDetails(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetails", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetails", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkScopesGetValidationDetails(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetails", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkScopesGetValidationDetails prepares the PrivateLinkScopesGetValidationDetails request.
func (c PrivateLinkScopesClient) preparerForPrivateLinkScopesGetValidationDetails(ctx context.Context, id PrivateLinkScopeId) (*http.Request, error) {
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

// responderForPrivateLinkScopesGetValidationDetails handles the response to the PrivateLinkScopesGetValidationDetails request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesClient) responderForPrivateLinkScopesGetValidationDetails(resp *http.Response) (result PrivateLinkScopesGetValidationDetailsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
