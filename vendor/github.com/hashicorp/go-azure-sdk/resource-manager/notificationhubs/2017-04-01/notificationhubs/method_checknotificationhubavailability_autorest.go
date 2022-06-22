package notificationhubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNotificationHubAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckAvailabilityResult
}

// CheckNotificationHubAvailability ...
func (c NotificationHubsClient) CheckNotificationHubAvailability(ctx context.Context, id NamespaceId, input CheckAvailabilityParameters) (result CheckNotificationHubAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckNotificationHubAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CheckNotificationHubAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CheckNotificationHubAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckNotificationHubAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "CheckNotificationHubAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckNotificationHubAvailability prepares the CheckNotificationHubAvailability request.
func (c NotificationHubsClient) preparerForCheckNotificationHubAvailability(ctx context.Context, id NamespaceId, input CheckAvailabilityParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkNotificationHubAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckNotificationHubAvailability handles the response to the CheckNotificationHubAvailability request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForCheckNotificationHubAvailability(resp *http.Response) (result CheckNotificationHubAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
