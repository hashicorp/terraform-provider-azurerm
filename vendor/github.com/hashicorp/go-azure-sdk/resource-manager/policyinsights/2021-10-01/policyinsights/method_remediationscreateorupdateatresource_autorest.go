package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsCreateOrUpdateAtResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCreateOrUpdateAtResource ...
func (c PolicyInsightsClient) RemediationsCreateOrUpdateAtResource(ctx context.Context, id ScopedRemediationId, input Remediation) (result RemediationsCreateOrUpdateAtResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsCreateOrUpdateAtResource(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResource", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResource", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCreateOrUpdateAtResource(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtResource", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCreateOrUpdateAtResource prepares the RemediationsCreateOrUpdateAtResource request.
func (c PolicyInsightsClient) preparerForRemediationsCreateOrUpdateAtResource(ctx context.Context, id ScopedRemediationId, input Remediation) (*http.Request, error) {
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

// responderForRemediationsCreateOrUpdateAtResource handles the response to the RemediationsCreateOrUpdateAtResource request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCreateOrUpdateAtResource(resp *http.Response) (result RemediationsCreateOrUpdateAtResourceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
