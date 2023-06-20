package workflowrunactions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRequestHistoriesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *RequestHistory
}

// WorkflowRunActionRequestHistoriesGet ...
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesGet(ctx context.Context, id RequestHistoryId) (result WorkflowRunActionRequestHistoriesGetOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRequestHistoriesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowRunActionRequestHistoriesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowRunActionRequestHistoriesGet prepares the WorkflowRunActionRequestHistoriesGet request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRequestHistoriesGet(ctx context.Context, id RequestHistoryId) (*http.Request, error) {
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

// responderForWorkflowRunActionRequestHistoriesGet handles the response to the WorkflowRunActionRequestHistoriesGet request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRequestHistoriesGet(resp *http.Response) (result WorkflowRunActionRequestHistoriesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
