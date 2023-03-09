package workflowrunactions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRepetitionsListExpressionTracesOperationResponse struct {
	HttpResponse *http.Response
	Model        *ExpressionTraces
}

// WorkflowRunActionRepetitionsListExpressionTraces ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsListExpressionTraces(ctx context.Context, id RepetitionId) (result WorkflowRunActionRepetitionsListExpressionTracesOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRepetitionsListExpressionTraces(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsListExpressionTraces", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsListExpressionTraces", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowRunActionRepetitionsListExpressionTraces(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsListExpressionTraces", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowRunActionRepetitionsListExpressionTraces prepares the WorkflowRunActionRepetitionsListExpressionTraces request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsListExpressionTraces(ctx context.Context, id RepetitionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listExpressionTraces", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkflowRunActionRepetitionsListExpressionTraces handles the response to the WorkflowRunActionRepetitionsListExpressionTraces request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRepetitionsListExpressionTraces(resp *http.Response) (result WorkflowRunActionRepetitionsListExpressionTracesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
