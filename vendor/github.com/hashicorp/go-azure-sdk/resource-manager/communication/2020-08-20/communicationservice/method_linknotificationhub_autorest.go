package communicationservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkNotificationHubOperationResponse struct {
	HttpResponse *http.Response
	Model        *LinkedNotificationHub
}

// LinkNotificationHub ...
func (c CommunicationServiceClient) LinkNotificationHub(ctx context.Context, id CommunicationServiceId, input LinkNotificationHubParameters) (result LinkNotificationHubOperationResponse, err error) {
	req, err := c.preparerForLinkNotificationHub(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "communicationservice.CommunicationServiceClient", "LinkNotificationHub", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "communicationservice.CommunicationServiceClient", "LinkNotificationHub", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLinkNotificationHub(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "communicationservice.CommunicationServiceClient", "LinkNotificationHub", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLinkNotificationHub prepares the LinkNotificationHub request.
func (c CommunicationServiceClient) preparerForLinkNotificationHub(ctx context.Context, id CommunicationServiceId, input LinkNotificationHubParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/linkNotificationHub", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLinkNotificationHub handles the response to the LinkNotificationHub request. The method always
// closes the http.Response Body.
func (c CommunicationServiceClient) responderForLinkNotificationHub(resp *http.Response) (result LinkNotificationHubOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
