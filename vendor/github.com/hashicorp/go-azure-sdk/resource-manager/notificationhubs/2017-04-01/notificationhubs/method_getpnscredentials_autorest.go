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

type GetPnsCredentialsOperationResponse struct {
	HttpResponse *http.Response
	Model        *PnsCredentialsResource
}

// GetPnsCredentials ...
func (c NotificationHubsClient) GetPnsCredentials(ctx context.Context, id NotificationHubId) (result GetPnsCredentialsOperationResponse, err error) {
	req, err := c.preparerForGetPnsCredentials(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "GetPnsCredentials", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "GetPnsCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetPnsCredentials(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NotificationHubsClient", "GetPnsCredentials", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetPnsCredentials prepares the GetPnsCredentials request.
func (c NotificationHubsClient) preparerForGetPnsCredentials(ctx context.Context, id NotificationHubId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/pnsCredentials", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetPnsCredentials handles the response to the GetPnsCredentials request. The method always
// closes the http.Response Body.
func (c NotificationHubsClient) responderForGetPnsCredentials(resp *http.Response) (result GetPnsCredentialsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
