package hybridconnections

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// GetAuthorizationRule ...
func (c HybridConnectionsClient) GetAuthorizationRule(ctx context.Context, id HybridConnectionAuthorizationRuleId) (result GetAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForGetAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "GetAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "GetAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "GetAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAuthorizationRule prepares the GetAuthorizationRule request.
func (c HybridConnectionsClient) preparerForGetAuthorizationRule(ctx context.Context, id HybridConnectionAuthorizationRuleId) (*http.Request, error) {
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

// responderForGetAuthorizationRule handles the response to the GetAuthorizationRule request. The method always
// closes the http.Response Body.
func (c HybridConnectionsClient) responderForGetAuthorizationRule(resp *http.Response) (result GetAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
