package policyassignments

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteByIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *PolicyAssignment
}

// DeleteById ...
func (c PolicyAssignmentsClient) DeleteById(ctx context.Context, id PolicyAssignmentIdId) (result DeleteByIdOperationResponse, err error) {
	req, err := c.preparerForDeleteById(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "DeleteById", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "DeleteById", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteById(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "DeleteById", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteById prepares the DeleteById request.
func (c PolicyAssignmentsClient) preparerForDeleteById(ctx context.Context, id PolicyAssignmentIdId) (*http.Request, error) {
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

// responderForDeleteById handles the response to the DeleteById request. The method always
// closes the http.Response Body.
func (c PolicyAssignmentsClient) responderForDeleteById(resp *http.Response) (result DeleteByIdOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
