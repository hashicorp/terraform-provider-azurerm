package policyinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsCreateOrUpdateAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCreateOrUpdateAtSubscription ...
func (c PolicyInsightsClient) RemediationsCreateOrUpdateAtSubscription(ctx context.Context, id RemediationId, input Remediation) (result RemediationsCreateOrUpdateAtSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsCreateOrUpdateAtSubscription(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCreateOrUpdateAtSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCreateOrUpdateAtSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCreateOrUpdateAtSubscription prepares the RemediationsCreateOrUpdateAtSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsCreateOrUpdateAtSubscription(ctx context.Context, id RemediationId, input Remediation) (*http.Request, error) {
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

// responderForRemediationsCreateOrUpdateAtSubscription handles the response to the RemediationsCreateOrUpdateAtSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCreateOrUpdateAtSubscription(resp *http.Response) (result RemediationsCreateOrUpdateAtSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
