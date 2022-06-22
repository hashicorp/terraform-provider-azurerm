package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsGetAtResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsGetAtResource ...
func (c PolicyInsightsClient) RemediationsGetAtResource(ctx context.Context, id ScopedRemediationId) (result RemediationsGetAtResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsGetAtResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResource", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResource", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsGetAtResource(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtResource", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsGetAtResource prepares the RemediationsGetAtResource request.
func (c PolicyInsightsClient) preparerForRemediationsGetAtResource(ctx context.Context, id ScopedRemediationId) (*http.Request, error) {
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

// responderForRemediationsGetAtResource handles the response to the RemediationsGetAtResource request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsGetAtResource(resp *http.Response) (result RemediationsGetAtResourceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
