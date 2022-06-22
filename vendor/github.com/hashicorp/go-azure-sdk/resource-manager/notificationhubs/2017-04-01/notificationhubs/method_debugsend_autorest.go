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

type DebugSendOperationResponse struct {
	HttpResponse *http.Response
	Model        *DebugSendResponse
}

// DebugSend ...
func (c NotificationHubsClient) DebugSend(ctx context.Context, id NotificationHubId, input interface{}) (result DebugSendOperationResponse, err error) {
	req, err := c.preparerForDebugSend(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DebugSend", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DebugSend", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDebugSend(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "DebugSend", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDebugSend prepares the DebugSend request.
func (c NotificationHubsClient) preparerForDebugSend(ctx context.Context, id NotificationHubId, input interface{}) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/debugsend", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDebugSend handles the response to the DebugSend request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForDebugSend(resp *http.Response) (result DebugSendOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
