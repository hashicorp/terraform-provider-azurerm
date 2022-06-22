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

type RemediationsCancelAtResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCancelAtResource ...
func (c PolicyInsightsClient) RemediationsCancelAtResource(ctx context.Context, id ScopedRemediationId) (result RemediationsCancelAtResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsCancelAtResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResource", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResource", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCancelAtResource(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResource", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCancelAtResource prepares the RemediationsCancelAtResource request.
func (c PolicyInsightsClient) preparerForRemediationsCancelAtResource(ctx context.Context, id ScopedRemediationId) (*http.Request, error) {
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

// responderForRemediationsCancelAtResource handles the response to the RemediationsCancelAtResource request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCancelAtResource(resp *http.Response) (result RemediationsCancelAtResourceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
