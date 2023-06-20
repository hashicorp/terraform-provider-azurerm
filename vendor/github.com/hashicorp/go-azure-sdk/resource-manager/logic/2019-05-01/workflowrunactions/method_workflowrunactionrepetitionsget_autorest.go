package workflowrunactions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRepetitionsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowRunActionRepetitionDefinition
}

// WorkflowRunActionRepetitionsGet ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsGet(ctx context.Context, id RepetitionId) (result WorkflowRunActionRepetitionsGetOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRepetitionsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowRunActionRepetitionsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowRunActionRepetitionsGet prepares the WorkflowRunActionRepetitionsGet request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsGet(ctx context.Context, id RepetitionId) (*http.Request, error) {
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

// responderForWorkflowRunActionRepetitionsGet handles the response to the WorkflowRunActionRepetitionsGet request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRepetitionsGet(resp *http.Response) (result WorkflowRunActionRepetitionsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
