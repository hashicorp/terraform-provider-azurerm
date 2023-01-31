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

type ListExpressionTracesOperationResponse struct {
	HttpResponse *http.Response
	Model        *ExpressionTraces
}

// ListExpressionTraces ...
func (c WorkflowRunActionsClient) ListExpressionTraces(ctx context.Context, id ActionId) (result ListExpressionTracesOperationResponse, err error) {
	req, err := c.preparerForListExpressionTraces(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "ListExpressionTraces", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "ListExpressionTraces", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListExpressionTraces(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "ListExpressionTraces", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListExpressionTraces prepares the ListExpressionTraces request.
func (c WorkflowRunActionsClient) preparerForListExpressionTraces(ctx context.Context, id ActionId) (*http.Request, error) {
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

// responderForListExpressionTraces handles the response to the ListExpressionTraces request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForListExpressionTraces(resp *http.Response) (result ListExpressionTracesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
