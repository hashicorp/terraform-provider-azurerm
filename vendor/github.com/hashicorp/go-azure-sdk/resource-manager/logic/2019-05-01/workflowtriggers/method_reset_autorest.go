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

type ResetOperationResponse struct {
	HttpResponse *http.Response
}

// Reset ...
func (c WorkflowTriggersClient) Reset(ctx context.Context, id TriggerId) (result ResetOperationResponse, err error) {
	req, err := c.preparerForReset(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Reset", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Reset", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForReset(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "Reset", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForReset prepares the Reset request.
func (c WorkflowTriggersClient) preparerForReset(ctx context.Context, id TriggerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reset", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForReset handles the response to the Reset request. The method always
// closes the http.Response Body.
func (c WorkflowTriggersClient) responderForReset(resp *http.Response) (result ResetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
