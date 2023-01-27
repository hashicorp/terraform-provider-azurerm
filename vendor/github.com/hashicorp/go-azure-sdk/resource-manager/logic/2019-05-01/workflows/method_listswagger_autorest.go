package workflows

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSwaggerOperationResponse struct {
	HttpResponse *http.Response
	Model        *interface{}
}

// ListSwagger ...
func (c WorkflowsClient) ListSwagger(ctx context.Context, id WorkflowId) (result ListSwaggerOperationResponse, err error) {
	req, err := c.preparerForListSwagger(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ListSwagger", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ListSwagger", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSwagger(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ListSwagger", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSwagger prepares the ListSwagger request.
func (c WorkflowsClient) preparerForListSwagger(ctx context.Context, id WorkflowId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listSwagger", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSwagger handles the response to the ListSwagger request. The method always
// closes the http.Response Body.
func (c WorkflowsClient) responderForListSwagger(resp *http.Response) (result ListSwaggerOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
