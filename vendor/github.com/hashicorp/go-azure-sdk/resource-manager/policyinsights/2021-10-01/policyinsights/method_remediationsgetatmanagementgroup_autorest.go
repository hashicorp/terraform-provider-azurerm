package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsGetAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsGetAtManagementGroup ...
func (c PolicyInsightsClient) RemediationsGetAtManagementGroup(ctx context.Context, id Providers2RemediationId) (result RemediationsGetAtManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsGetAtManagementGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtManagementGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtManagementGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsGetAtManagementGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtManagementGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsGetAtManagementGroup prepares the RemediationsGetAtManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsGetAtManagementGroup(ctx context.Context, id Providers2RemediationId) (*http.Request, error) {
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

// responderForRemediationsGetAtManagementGroup handles the response to the RemediationsGetAtManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsGetAtManagementGroup(resp *http.Response) (result RemediationsGetAtManagementGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
