package authorizationrulesnamespaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespacesDeleteAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
}

// NamespacesDeleteAuthorizationRule ...
func (c AuthorizationRulesNamespacesClient) NamespacesDeleteAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (result NamespacesDeleteAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForNamespacesDeleteAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesDeleteAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesDeleteAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesDeleteAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesDeleteAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesDeleteAuthorizationRule prepares the NamespacesDeleteAuthorizationRule request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesDeleteAuthorizationRule(ctx context.Context, id AuthorizationRuleId) (*http.Request, error) {
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

// responderForNamespacesDeleteAuthorizationRule handles the response to the NamespacesDeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesDeleteAuthorizationRule(resp *http.Response) (result NamespacesDeleteAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
