package authorizationruleseventhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsCreateOrUpdateAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *AuthorizationRule
}

// EventHubsCreateOrUpdateAuthorizationRule ...
func (c AuthorizationRulesEventHubsClient) EventHubsCreateOrUpdateAuthorizationRule(ctx context.Context, id EventhubAuthorizationRuleId, input AuthorizationRule) (result EventHubsCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForEventHubsCreateOrUpdateAuthorizationRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsCreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEventHubsCreateOrUpdateAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsCreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEventHubsCreateOrUpdateAuthorizationRule prepares the EventHubsCreateOrUpdateAuthorizationRule request.
func (c AuthorizationRulesEventHubsClient) preparerForEventHubsCreateOrUpdateAuthorizationRule(ctx context.Context, id EventhubAuthorizationRuleId, input AuthorizationRule) (*http.Request, error) {
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

// responderForEventHubsCreateOrUpdateAuthorizationRule handles the response to the EventHubsCreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesEventHubsClient) responderForEventHubsCreateOrUpdateAuthorizationRule(resp *http.Response) (result EventHubsCreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
