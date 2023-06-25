package policyassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetByIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *PolicyAssignment
}

// GetById ...
func (c PolicyAssignmentsClient) GetById(ctx context.Context, id PolicyAssignmentIdId) (result GetByIdOperationResponse, err error) {
	req, err := c.preparerForGetById(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "GetById", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "GetById", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetById(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "GetById", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetById prepares the GetById request.
func (c PolicyAssignmentsClient) preparerForGetById(ctx context.Context, id PolicyAssignmentIdId) (*http.Request, error) {
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

// responderForGetById handles the response to the GetById request. The method always
// closes the http.Response Body.
func (c PolicyAssignmentsClient) responderForGetById(resp *http.Response) (result GetByIdOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
