package workflowrunactions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopeRepetitionsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowRunActionRepetitionDefinition
}

// CopeRepetitionsGet ...
func (c WorkflowRunActionsClient) CopeRepetitionsGet(ctx context.Context, id ScopeRepetitionId) (result CopeRepetitionsGetOperationResponse, err error) {
	req, err := c.preparerForCopeRepetitionsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCopeRepetitionsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCopeRepetitionsGet prepares the CopeRepetitionsGet request.
func (c WorkflowRunActionsClient) preparerForCopeRepetitionsGet(ctx context.Context, id ScopeRepetitionId) (*http.Request, error) {
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

// responderForCopeRepetitionsGet handles the response to the CopeRepetitionsGet request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForCopeRepetitionsGet(resp *http.Response) (result CopeRepetitionsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
