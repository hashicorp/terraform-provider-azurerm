package notificationhubs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchOperationResponse struct {
	HttpResponse *http.Response
	Model        *NotificationHubResource
}

// Patch ...
func (c NotificationHubsClient) Patch(ctx context.Context, id NotificationHubId, input NotificationHubPatchParameters) (result PatchOperationResponse, err error) {
	req, err := c.preparerForPatch(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "Patch", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "Patch", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatch(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "Patch", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatch prepares the Patch request.
func (c NotificationHubsClient) preparerForPatch(ctx context.Context, id NotificationHubId, input NotificationHubPatchParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatch handles the response to the Patch request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForPatch(resp *http.Response) (result PatchOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
