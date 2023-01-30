package roleeligibilityschedulerequests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CancelOperationResponse struct {
	HttpResponse *http.Response
}

// Cancel ...
func (c RoleEligibilityScheduleRequestsClient) Cancel(ctx context.Context, id ScopedRoleEligibilityScheduleRequestId) (result CancelOperationResponse, err error) {
	req, err := c.preparerForCancel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient", "Cancel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient", "Cancel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCancel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleeligibilityschedulerequests.RoleEligibilityScheduleRequestsClient", "Cancel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCancel prepares the Cancel request.
func (c RoleEligibilityScheduleRequestsClient) preparerForCancel(ctx context.Context, id ScopedRoleEligibilityScheduleRequestId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cancel", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCancel handles the response to the Cancel request. The method always
// closes the http.Response Body.
func (c RoleEligibilityScheduleRequestsClient) responderForCancel(resp *http.Response) (result CancelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
