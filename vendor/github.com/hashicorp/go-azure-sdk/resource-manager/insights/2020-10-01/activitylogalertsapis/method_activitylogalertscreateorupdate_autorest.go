package activitylogalertsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogAlertsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActivityLogAlertResource
}

// ActivityLogAlertsCreateOrUpdate ...
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsCreateOrUpdate(ctx context.Context, id ActivityLogAlertId, input ActivityLogAlertResource) (result ActivityLogAlertsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForActivityLogAlertsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActivityLogAlertsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActivityLogAlertsCreateOrUpdate prepares the ActivityLogAlertsCreateOrUpdate request.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsCreateOrUpdate(ctx context.Context, id ActivityLogAlertId, input ActivityLogAlertResource) (*http.Request, error) {
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

// responderForActivityLogAlertsCreateOrUpdate handles the response to the ActivityLogAlertsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ActivityLogAlertsAPIsClient) responderForActivityLogAlertsCreateOrUpdate(resp *http.Response) (result ActivityLogAlertsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
