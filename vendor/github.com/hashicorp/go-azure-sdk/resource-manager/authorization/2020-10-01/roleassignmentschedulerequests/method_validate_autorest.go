package roleassignmentschedulerequests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateOperationResponse struct {
	HttpResponse *http.Response
	Model        *RoleAssignmentScheduleRequest
}

// Validate ...
func (c RoleAssignmentScheduleRequestsClient) Validate(ctx context.Context, id ScopedRoleAssignmentScheduleRequestId, input RoleAssignmentScheduleRequest) (result ValidateOperationResponse, err error) {
	req, err := c.preparerForValidate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient", "Validate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient", "Validate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "roleassignmentschedulerequests.RoleAssignmentScheduleRequestsClient", "Validate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidate prepares the Validate request.
func (c RoleAssignmentScheduleRequestsClient) preparerForValidate(ctx context.Context, id ScopedRoleAssignmentScheduleRequestId, input RoleAssignmentScheduleRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForValidate handles the response to the Validate request. The method always
// closes the http.Response Body.
func (c RoleAssignmentScheduleRequestsClient) responderForValidate(resp *http.Response) (result ValidateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
