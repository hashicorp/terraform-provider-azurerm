package queuesauthorizationrule

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesDeleteAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
}

// QueuesDeleteAuthorizationRule ...
func (c QueuesAuthorizationRuleClient) QueuesDeleteAuthorizationRule(ctx context.Context, id QueueAuthorizationRuleId) (result QueuesDeleteAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForQueuesDeleteAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesDeleteAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesDeleteAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueuesDeleteAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesDeleteAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueuesDeleteAuthorizationRule prepares the QueuesDeleteAuthorizationRule request.
func (c QueuesAuthorizationRuleClient) preparerForQueuesDeleteAuthorizationRule(ctx context.Context, id QueueAuthorizationRuleId) (*http.Request, error) {
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

// responderForQueuesDeleteAuthorizationRule handles the response to the QueuesDeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (c QueuesAuthorizationRuleClient) responderForQueuesDeleteAuthorizationRule(resp *http.Response) (result QueuesDeleteAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
