package activitylogalertsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogAlertsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActivityLogAlertResource
}

// ActivityLogAlertsGet ...
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsGet(ctx context.Context, id ActivityLogAlertId) (result ActivityLogAlertsGetOperationResponse, err error) {
	req, err := c.preparerForActivityLogAlertsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActivityLogAlertsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActivityLogAlertsGet prepares the ActivityLogAlertsGet request.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsGet(ctx context.Context, id ActivityLogAlertId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForActivityLogAlertsGet handles the response to the ActivityLogAlertsGet request. The method always
// closes the http.Response Body.
func (c ActivityLogAlertsAPIsClient) responderForActivityLogAlertsGet(resp *http.Response) (result ActivityLogAlertsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
