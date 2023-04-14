package policyassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateByIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *PolicyAssignment
}

// CreateById ...
func (c PolicyAssignmentsClient) CreateById(ctx context.Context, id PolicyAssignmentIdId, input PolicyAssignment) (result CreateByIdOperationResponse, err error) {
	req, err := c.preparerForCreateById(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "CreateById", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "CreateById", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateById(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "CreateById", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateById prepares the CreateById request.
func (c PolicyAssignmentsClient) preparerForCreateById(ctx context.Context, id PolicyAssignmentIdId, input PolicyAssignment) (*http.Request, error) {
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

// responderForCreateById handles the response to the CreateById request. The method always
// closes the http.Response Body.
func (c PolicyAssignmentsClient) responderForCreateById(resp *http.Response) (result CreateByIdOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
