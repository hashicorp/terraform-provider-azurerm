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

type EnableOperationResponse struct {
	HttpResponse *http.Response
}

// Enable ...
func (c WorkflowsClient) Enable(ctx context.Context, id WorkflowId) (result EnableOperationResponse, err error) {
	req, err := c.preparerForEnable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "Enable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "Enable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForEnable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "Enable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForEnable prepares the Enable request.
func (c WorkflowsClient) preparerForEnable(ctx context.Context, id WorkflowId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForEnable handles the response to the Enable request. The method always
// closes the http.Response Body.
func (c WorkflowsClient) responderForEnable(resp *http.Response) (result EnableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
