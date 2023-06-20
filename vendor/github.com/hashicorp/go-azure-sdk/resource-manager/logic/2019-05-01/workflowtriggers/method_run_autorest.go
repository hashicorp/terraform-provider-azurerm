package workflowtriggers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunOperationResponse struct {
	HttpResponse *http.Response
}

// Run ...
func (c WorkflowTriggersClient) Run(ctx context.Context, id TriggerId) (result RunOperationResponse, err error) {
	req, err := c.preparerForRun(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Run", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Run", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRun(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Run", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRun prepares the Run request.
func (c WorkflowTriggersClient) preparerForRun(ctx context.Context, id TriggerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/run", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRun handles the response to the Run request. The method always
// closes the http.Response Body.
func (c WorkflowTriggersClient) responderForRun(resp *http.Response) (result RunOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
