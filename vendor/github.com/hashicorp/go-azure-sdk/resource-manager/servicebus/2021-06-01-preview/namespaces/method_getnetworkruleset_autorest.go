package namespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetNetworkRuleSetOperationResponse struct {
	HttpResponse *http.Response
	Model        *NetworkRuleSet
}

// GetNetworkRuleSet ...
func (c NamespacesClient) GetNetworkRuleSet(ctx context.Context, id NamespaceId) (result GetNetworkRuleSetOperationResponse, err error) {
	req, err := c.preparerForGetNetworkRuleSet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "GetNetworkRuleSet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "GetNetworkRuleSet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetNetworkRuleSet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "GetNetworkRuleSet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetNetworkRuleSet prepares the GetNetworkRuleSet request.
func (c NamespacesClient) preparerForGetNetworkRuleSet(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetNetworkRuleSet handles the response to the GetNetworkRuleSet request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForGetNetworkRuleSet(resp *http.Response) (result GetNetworkRuleSetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
