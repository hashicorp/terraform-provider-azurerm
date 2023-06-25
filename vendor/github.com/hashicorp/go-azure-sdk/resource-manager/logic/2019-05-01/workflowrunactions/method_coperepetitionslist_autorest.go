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

type CopeRepetitionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowRunActionRepetitionDefinitionCollection
}

// CopeRepetitionsList ...
func (c WorkflowRunActionsClient) CopeRepetitionsList(ctx context.Context, id ActionId) (result CopeRepetitionsListOperationResponse, err error) {
	req, err := c.preparerForCopeRepetitionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCopeRepetitionsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "CopeRepetitionsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCopeRepetitionsList prepares the CopeRepetitionsList request.
func (c WorkflowRunActionsClient) preparerForCopeRepetitionsList(ctx context.Context, id ActionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/scopeRepetitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCopeRepetitionsList handles the response to the CopeRepetitionsList request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForCopeRepetitionsList(resp *http.Response) (result CopeRepetitionsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
