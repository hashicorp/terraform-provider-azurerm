package runbookdraft

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UndoEditOperationResponse struct {
	HttpResponse *http.Response
}

// UndoEdit ...
func (c RunbookDraftClient) UndoEdit(ctx context.Context, id RunbookId) (result UndoEditOperationResponse, err error) {
	req, err := c.preparerForUndoEdit(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbookdraft.RunbookDraftClient", "UndoEdit", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbookdraft.RunbookDraftClient", "UndoEdit", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUndoEdit(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "runbookdraft.RunbookDraftClient", "UndoEdit", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUndoEdit prepares the UndoEdit request.
func (c RunbookDraftClient) preparerForUndoEdit(ctx context.Context, id RunbookId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/draft/undoEdit", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUndoEdit handles the response to the UndoEdit request. The method always
// closes the http.Response Body.
func (c RunbookDraftClient) responderForUndoEdit(resp *http.Response) (result UndoEditOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
