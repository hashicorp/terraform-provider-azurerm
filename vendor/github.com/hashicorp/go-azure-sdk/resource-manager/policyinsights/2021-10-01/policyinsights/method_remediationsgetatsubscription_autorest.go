package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsGetAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsGetAtSubscription ...
func (c PolicyInsightsClient) RemediationsGetAtSubscription(ctx context.Context, id RemediationId) (result RemediationsGetAtSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsGetAtSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsGetAtSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsGetAtSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsGetAtSubscription prepares the RemediationsGetAtSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsGetAtSubscription(ctx context.Context, id RemediationId) (*http.Request, error) {
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

// responderForRemediationsGetAtSubscription handles the response to the RemediationsGetAtSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsGetAtSubscription(resp *http.Response) (result RemediationsGetAtSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
