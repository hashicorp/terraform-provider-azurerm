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

type WorkflowRunActionRepetitionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowRunActionRepetitionDefinitionCollection
}

// WorkflowRunActionRepetitionsList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsList(ctx context.Context, id ActionId) (result WorkflowRunActionRepetitionsListOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRepetitionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowRunActionRepetitionsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowRunActionRepetitionsList prepares the WorkflowRunActionRepetitionsList request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsList(ctx context.Context, id ActionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/repetitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkflowRunActionRepetitionsList handles the response to the WorkflowRunActionRepetitionsList request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRepetitionsList(resp *http.Response) (result WorkflowRunActionRepetitionsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
