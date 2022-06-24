package notificationhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteAuthorizationRuleOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteAuthorizationRule ...
func (c NotificationHubsClient) DeleteAuthorizationRule(ctx context.Context, id NotificationHubAuthorizationRuleId) (result DeleteAuthorizationRuleOperationResponse, err error) {
	req, err := c.preparerForDeleteAuthorizationRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DeleteAuthorizationRule", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DeleteAuthorizationRule", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteAuthorizationRule(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DeleteAuthorizationRule", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteAuthorizationRule prepares the DeleteAuthorizationRule request.
func (c NotificationHubsClient) preparerForDeleteAuthorizationRule(ctx context.Context, id NotificationHubAuthorizationRuleId) (*http.Request, error) {
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

// responderForDeleteAuthorizationRule handles the response to the DeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForDeleteAuthorizationRule(resp *http.Response) (result DeleteAuthorizationRuleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
