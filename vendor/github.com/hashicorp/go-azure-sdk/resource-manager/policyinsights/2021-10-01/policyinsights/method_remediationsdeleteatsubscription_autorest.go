package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsDeleteAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsDeleteAtSubscription ...
func (c PolicyInsightsClient) RemediationsDeleteAtSubscription(ctx context.Context, id RemediationId) (result RemediationsDeleteAtSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsDeleteAtSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsDeleteAtSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsDeleteAtSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsDeleteAtSubscription prepares the RemediationsDeleteAtSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsDeleteAtSubscription(ctx context.Context, id RemediationId) (*http.Request, error) {
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

// responderForRemediationsDeleteAtSubscription handles the response to the RemediationsDeleteAtSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsDeleteAtSubscription(resp *http.Response) (result RemediationsDeleteAtSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
