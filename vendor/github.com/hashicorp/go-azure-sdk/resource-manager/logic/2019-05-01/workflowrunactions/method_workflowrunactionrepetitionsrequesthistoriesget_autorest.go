package workflowrunactions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRepetitionsRequestHistoriesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *RequestHistory
}

// WorkflowRunActionRepetitionsRequestHistoriesGet ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesGet(ctx context.Context, id RepetitionRequestHistoryId) (result WorkflowRunActionRepetitionsRequestHistoriesGetOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRepetitionsRequestHistoriesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowRunActionRepetitionsRequestHistoriesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowRunActionRepetitionsRequestHistoriesGet prepares the WorkflowRunActionRepetitionsRequestHistoriesGet request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsRequestHistoriesGet(ctx context.Context, id RepetitionRequestHistoryId) (*http.Request, error) {
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

// responderForWorkflowRunActionRepetitionsRequestHistoriesGet handles the response to the WorkflowRunActionRepetitionsRequestHistoriesGet request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRepetitionsRequestHistoriesGet(resp *http.Response) (result WorkflowRunActionRepetitionsRequestHistoriesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
