package remediations

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationsDeleteAtResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *Remediation
}

// RemediationsDeleteAtResource ...
func (c RemediationsClient) RemediationsDeleteAtResource(ctx context.Context, id ScopedRemediationId) (result RemediationsDeleteAtResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsDeleteAtResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "remediations.RemediationsClient", "RemediationsDeleteAtResource", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "remediations.RemediationsClient", "RemediationsDeleteAtResource", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemediationsDeleteAtResource(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "remediations.RemediationsClient", "RemediationsDeleteAtResource", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemediationsDeleteAtResource prepares the RemediationsDeleteAtResource request.
func (c RemediationsClient) preparerForRemediationsDeleteAtResource(ctx context.Context, id ScopedRemediationId) (*http.Request, error) {
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

// responderForRemediationsDeleteAtResource handles the response to the RemediationsDeleteAtResource request. The method always
// closes the http.Response Body.
func (c RemediationsClient) responderForRemediationsDeleteAtResource(resp *http.Response) (result RemediationsDeleteAtResourceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
