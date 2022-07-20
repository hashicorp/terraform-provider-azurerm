package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsCreateOrUpdateAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCreateOrUpdateAtManagementGroup ...
func (c PolicyInsightsClient) RemediationsCreateOrUpdateAtManagementGroup(ctx context.Context, id Providers2RemediationId, input Remediation) (result RemediationsCreateOrUpdateAtManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsCreateOrUpdateAtManagementGroup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtManagementGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtManagementGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCreateOrUpdateAtManagementGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtManagementGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCreateOrUpdateAtManagementGroup prepares the RemediationsCreateOrUpdateAtManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsCreateOrUpdateAtManagementGroup(ctx context.Context, id Providers2RemediationId, input Remediation) (*http.Request, error) {
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

// responderForRemediationsCreateOrUpdateAtManagementGroup handles the response to the RemediationsCreateOrUpdateAtManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCreateOrUpdateAtManagementGroup(resp *http.Response) (result RemediationsCreateOrUpdateAtManagementGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
