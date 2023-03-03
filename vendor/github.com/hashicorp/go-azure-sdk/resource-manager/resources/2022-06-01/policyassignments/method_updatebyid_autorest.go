package policyassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateByIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *PolicyAssignment
}

// UpdateById ...
func (c PolicyAssignmentsClient) UpdateById(ctx context.Context, id PolicyAssignmentIdId, input PolicyAssignmentUpdate) (result UpdateByIdOperationResponse, err error) {
	req, err := c.preparerForUpdateById(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "UpdateById", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "UpdateById", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdateById(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "UpdateById", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdateById prepares the UpdateById request.
func (c PolicyAssignmentsClient) preparerForUpdateById(ctx context.Context, id PolicyAssignmentIdId, input PolicyAssignmentUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdateById handles the response to the UpdateById request. The method always
// closes the http.Response Body.
func (c PolicyAssignmentsClient) responderForUpdateById(resp *http.Response) (result UpdateByIdOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
