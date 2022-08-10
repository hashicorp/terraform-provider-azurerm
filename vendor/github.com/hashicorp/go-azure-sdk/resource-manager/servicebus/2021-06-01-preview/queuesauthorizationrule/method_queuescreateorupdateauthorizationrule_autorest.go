package queuesauthorizationrule

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesCreateOrUpdateAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *SBAuthorizationRule
}

// QueuesCreateOrUpdateAuthorizationRule ...
func (c QueuesAuthorizationRuleClient) QueuesCreateOrUpdateAuthorizationRule(ctx context.Context, id QueueAuthorizationRuleId, input SBAuthorizationRule) (result QueuesCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForQueuesCreateOrUpdateAuthorizationRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesCreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueuesCreateOrUpdateAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueuesCreateOrUpdateAuthorizationRule prepares the QueuesCreateOrUpdateAuthorizationRule request.
func (c QueuesAuthorizationRuleClient) preparerForQueuesCreateOrUpdateAuthorizationRule(ctx context.Context, id QueueAuthorizationRuleId, input SBAuthorizationRule) (*http.Request, error) {
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

// responderForQueuesCreateOrUpdateAuthorizationRule handles the response to the QueuesCreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (c QueuesAuthorizationRuleClient) responderForQueuesCreateOrUpdateAuthorizationRule(resp *http.Response) (result QueuesCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
