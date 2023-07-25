package activitylogalertsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogAlertsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ActivityLogAlertsDelete ...
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsDelete(ctx context.Context, id ActivityLogAlertId) (result ActivityLogAlertsDeleteOperationResponse, err error) {
	req, err := c.preparerForActivityLogAlertsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActivityLogAlertsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActivityLogAlertsDelete prepares the ActivityLogAlertsDelete request.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsDelete(ctx context.Context, id ActivityLogAlertId) (*http.Request, error) {
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

// responderForActivityLogAlertsDelete handles the response to the ActivityLogAlertsDelete request. The method always
// closes the http.Response Body.
func (c ActivityLogAlertsAPIsClient) responderForActivityLogAlertsDelete(resp *http.Response) (result ActivityLogAlertsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
