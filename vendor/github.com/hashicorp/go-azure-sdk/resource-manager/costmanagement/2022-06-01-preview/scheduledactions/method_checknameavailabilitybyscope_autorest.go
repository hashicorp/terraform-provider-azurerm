package scheduledactions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityByScopeOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResponse
}

// CheckNameAvailabilityByScope ...
func (c ScheduledActionsClient) CheckNameAvailabilityByScope(ctx context.Context, id commonids.ScopeId, input CheckNameAvailabilityRequest) (result CheckNameAvailabilityByScopeOperationResponse, err error) {
	req, err := c.preparerForCheckNameAvailabilityByScope(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailabilityByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailabilityByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckNameAvailabilityByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "CheckNameAvailabilityByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckNameAvailabilityByScope prepares the CheckNameAvailabilityByScope request.
func (c ScheduledActionsClient) preparerForCheckNameAvailabilityByScope(ctx context.Context, id commonids.ScopeId, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CostManagement/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckNameAvailabilityByScope handles the response to the CheckNameAvailabilityByScope request. The method always
// closes the http.Response Body.
func (c ScheduledActionsClient) responderForCheckNameAvailabilityByScope(resp *http.Response) (result CheckNameAvailabilityByScopeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
