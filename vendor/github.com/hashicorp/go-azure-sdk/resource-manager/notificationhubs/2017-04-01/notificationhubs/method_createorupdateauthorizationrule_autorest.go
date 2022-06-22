package notificationhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrUpdateAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedAccessAuthorizationRuleResource
}

// CreateOrUpdateAuthorizationRule ...
func (c NotificationHubsClient) CreateOrUpdateAuthorizationRule(ctx context.Context, id NotificationHubAuthorizationRuleId, input SharedAccessAuthorizationRuleCreateOrUpdateParameters) (result CreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateAuthorizationRule(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CreateOrUpdateAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateAuthorizationRule prepares the CreateOrUpdateAuthorizationRule request.
func (c NotificationHubsClient) preparerForCreateOrUpdateAuthorizationRule(ctx context.Context, id NotificationHubAuthorizationRuleId, input SharedAccessAuthorizationRuleCreateOrUpdateParameters) (*http.Request, error) {
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

// responderForCreateOrUpdateAuthorizationRule handles the response to the CreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForCreateOrUpdateAuthorizationRule(resp *http.Response) (result CreateOrUpdateAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
