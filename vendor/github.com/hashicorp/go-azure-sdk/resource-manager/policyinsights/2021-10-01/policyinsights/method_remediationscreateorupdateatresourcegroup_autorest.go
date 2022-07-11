package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsCreateOrUpdateAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCreateOrUpdateAtResourceGroup ...
func (c PolicyInsightsClient) RemediationsCreateOrUpdateAtResourceGroup(ctx context.Context, id ProviderRemediationId, input Remediation) (result RemediationsCreateOrUpdateAtResourceGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsCreateOrUpdateAtResourceGroup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCreateOrUpdateAtResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCreateOrUpdateAtResourceGroup prepares the RemediationsCreateOrUpdateAtResourceGroup request.
func (c PolicyInsightsClient) preparerForRemediationsCreateOrUpdateAtResourceGroup(ctx context.Context, id ProviderRemediationId, input Remediation) (*http.Request, error) {
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

// responderForRemediationsCreateOrUpdateAtResourceGroup handles the response to the RemediationsCreateOrUpdateAtResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCreateOrUpdateAtResourceGroup(resp *http.Response) (result RemediationsCreateOrUpdateAtResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
