package scheduledactions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResponse
}

// CheckNameAvailability ...
func (c ScheduledActionsClient) CheckNameAvailability(ctx context.Context, input CheckNameAvailabilityRequest) (result CheckNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckNameAvailability(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckNameAvailability prepares the CheckNameAvailability request.
func (c ScheduledActionsClient) preparerForCheckNameAvailability(ctx context.Context, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.CostManagement/checkNameAvailability"),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckNameAvailability handles the response to the CheckNameAvailability request. The method always
// closes the http.Response Body.
func (c ScheduledActionsClient) responderForCheckNameAvailability(resp *http.Response) (result CheckNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
