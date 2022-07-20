package topicsauthorizationrule

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsCreateOrUpdateAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *SBAuthorizationRule
}

// TopicsCreateOrUpdateAuthorizationRule ...
func (c TopicsAuthorizationRuleClient) TopicsCreateOrUpdateAuthorizationRule(ctx context.Context, id TopicAuthorizationRuleId, input SBAuthorizationRule) (result TopicsCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForTopicsCreateOrUpdateAuthorizationRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsCreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTopicsCreateOrUpdateAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTopicsCreateOrUpdateAuthorizationRule prepares the TopicsCreateOrUpdateAuthorizationRule request.
func (c TopicsAuthorizationRuleClient) preparerForTopicsCreateOrUpdateAuthorizationRule(ctx context.Context, id TopicAuthorizationRuleId, input SBAuthorizationRule) (*http.Request, error) {
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

// responderForTopicsCreateOrUpdateAuthorizationRule handles the response to the TopicsCreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (c TopicsAuthorizationRuleClient) responderForTopicsCreateOrUpdateAuthorizationRule(resp *http.Response) (result TopicsCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
