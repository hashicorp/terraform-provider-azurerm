package authorizationruleseventhubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubsRegenerateKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// EventHubsRegenerateKeys ...
func (c AuthorizationRulesEventHubsClient) EventHubsRegenerateKeys(ctx context.Context, id EventhubAuthorizationRuleId, input RegenerateAccessKeyParameters) (result EventHubsRegenerateKeysOperationResponse, err error) {
	req, err := c.preparerForEventHubsRegenerateKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsRegenerateKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsRegenerateKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEventHubsRegenerateKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsRegenerateKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEventHubsRegenerateKeys prepares the EventHubsRegenerateKeys request.
func (c AuthorizationRulesEventHubsClient) preparerForEventHubsRegenerateKeys(ctx context.Context, id EventhubAuthorizationRuleId, input RegenerateAccessKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEventHubsRegenerateKeys handles the response to the EventHubsRegenerateKeys request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesEventHubsClient) responderForEventHubsRegenerateKeys(resp *http.Response) (result EventHubsRegenerateKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
