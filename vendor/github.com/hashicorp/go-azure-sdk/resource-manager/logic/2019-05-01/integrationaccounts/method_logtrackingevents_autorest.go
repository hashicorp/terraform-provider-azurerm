package integrationaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogTrackingEventsOperationResponse struct {
	HttpResponse *http.Response
}

// LogTrackingEvents ...
func (c IntegrationAccountsClient) LogTrackingEvents(ctx context.Context, id IntegrationAccountId, input TrackingEventsDefinition) (result LogTrackingEventsOperationResponse, err error) {
	req, err := c.preparerForLogTrackingEvents(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "LogTrackingEvents", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "LogTrackingEvents", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLogTrackingEvents(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "LogTrackingEvents", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLogTrackingEvents prepares the LogTrackingEvents request.
func (c IntegrationAccountsClient) preparerForLogTrackingEvents(ctx context.Context, id IntegrationAccountId, input TrackingEventsDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/logTrackingEvents", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLogTrackingEvents handles the response to the LogTrackingEvents request. The method always
// closes the http.Response Body.
func (c IntegrationAccountsClient) responderForLogTrackingEvents(resp *http.Response) (result LogTrackingEventsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
