package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsGetAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsGetAtResourceGroup ...
func (c PolicyInsightsClient) RemediationsGetAtResourceGroup(ctx context.Context, id ProviderRemediationId) (result RemediationsGetAtResourceGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsGetAtResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsGetAtResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsGetAtResourceGroup prepares the RemediationsGetAtResourceGroup request.
func (c PolicyInsightsClient) preparerForRemediationsGetAtResourceGroup(ctx context.Context, id ProviderRemediationId) (*http.Request, error) {
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

// responderForRemediationsGetAtResourceGroup handles the response to the RemediationsGetAtResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsGetAtResourceGroup(resp *http.Response) (result RemediationsGetAtResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
