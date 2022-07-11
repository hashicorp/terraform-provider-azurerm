package policyinsights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsCancelAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCancelAtSubscription ...
func (c PolicyInsightsClient) RemediationsCancelAtSubscription(ctx context.Context, id RemediationId) (result RemediationsCancelAtSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsCancelAtSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCancelAtSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCancelAtSubscription prepares the RemediationsCancelAtSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsCancelAtSubscription(ctx context.Context, id RemediationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cancel", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemediationsCancelAtSubscription handles the response to the RemediationsCancelAtSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCancelAtSubscription(resp *http.Response) (result RemediationsCancelAtSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
