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

type CreateOrUpdateNetworkRuleSetOperationResponse struct {
	HttpResponse *http.Response
	Model        *NetworkRuleSet
}

// CreateOrUpdateNetworkRuleSet ...
func (c NamespacesClient) CreateOrUpdateNetworkRuleSet(ctx context.Context, id NamespaceId, input NetworkRuleSet) (result CreateOrUpdateNetworkRuleSetOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateNetworkRuleSet(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CreateOrUpdateNetworkRuleSet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CreateOrUpdateNetworkRuleSet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateNetworkRuleSet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CreateOrUpdateNetworkRuleSet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateNetworkRuleSet prepares the CreateOrUpdateNetworkRuleSet request.
func (c NamespacesClient) preparerForCreateOrUpdateNetworkRuleSet(ctx context.Context, id NamespaceId, input NetworkRuleSet) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreateOrUpdateNetworkRuleSet handles the response to the CreateOrUpdateNetworkRuleSet request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForCreateOrUpdateNetworkRuleSet(resp *http.Response) (result CreateOrUpdateNetworkRuleSetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
