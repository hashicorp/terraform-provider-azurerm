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

type RemediationsCancelAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsCancelAtResourceGroup ...
func (c PolicyInsightsClient) RemediationsCancelAtResourceGroup(ctx context.Context, id ProviderRemediationId) (result RemediationsCancelAtResourceGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsCancelAtResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsCancelAtResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsCancelAtResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsCancelAtResourceGroup prepares the RemediationsCancelAtResourceGroup request.
func (c PolicyInsightsClient) preparerForRemediationsCancelAtResourceGroup(ctx context.Context, id ProviderRemediationId) (*http.Request, error) {
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

// responderForRemediationsCancelAtResourceGroup handles the response to the RemediationsCancelAtResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsCancelAtResourceGroup(resp *http.Response) (result RemediationsCancelAtResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
