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

type RemediationsCancelAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCancelAtManagementGroup ...
func (c PolicyInsightsClient) RemediationsCancelAtManagementGroup(ctx context.Context, id Providers2RemediationId) (result RemediationsCancelAtManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsCancelAtManagementGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtManagementGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtManagementGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCancelAtManagementGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtManagementGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCancelAtManagementGroup prepares the RemediationsCancelAtManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsCancelAtManagementGroup(ctx context.Context, id Providers2RemediationId) (*http.Request, error) {
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

// responderForRemediationsCancelAtManagementGroup handles the response to the RemediationsCancelAtManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCancelAtManagementGroup(resp *http.Response) (result RemediationsCancelAtManagementGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
