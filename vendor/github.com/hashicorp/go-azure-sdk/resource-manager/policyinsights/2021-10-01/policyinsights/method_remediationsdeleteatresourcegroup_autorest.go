package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsDeleteAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsDeleteAtResourceGroup ...
func (c PolicyInsightsClient) RemediationsDeleteAtResourceGroup(ctx context.Context, id ProviderRemediationId) (result RemediationsDeleteAtResourceGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsDeleteAtResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsDeleteAtResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsDeleteAtResourceGroup prepares the RemediationsDeleteAtResourceGroup request.
func (c PolicyInsightsClient) preparerForRemediationsDeleteAtResourceGroup(ctx context.Context, id ProviderRemediationId) (*http.Request, error) {
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

// responderForRemediationsDeleteAtResourceGroup handles the response to the RemediationsDeleteAtResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsDeleteAtResourceGroup(resp *http.Response) (result RemediationsDeleteAtResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
