package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsDeleteAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsDeleteAtManagementGroup ...
func (c PolicyInsightsClient) RemediationsDeleteAtManagementGroup(ctx context.Context, id Providers2RemediationId) (result RemediationsDeleteAtManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsDeleteAtManagementGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtManagementGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtManagementGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsDeleteAtManagementGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtManagementGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsDeleteAtManagementGroup prepares the RemediationsDeleteAtManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsDeleteAtManagementGroup(ctx context.Context, id Providers2RemediationId) (*http.Request, error) {
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

// responderForRemediationsDeleteAtManagementGroup handles the response to the RemediationsDeleteAtManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsDeleteAtManagementGroup(resp *http.Response) (result RemediationsDeleteAtManagementGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
