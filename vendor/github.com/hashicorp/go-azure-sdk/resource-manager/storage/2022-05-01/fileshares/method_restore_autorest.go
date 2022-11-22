package fileshares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreOperationResponse struct {
	HttpResponse *http.Response
}

// Restore ...
func (c FileSharesClient) Restore(ctx context.Context, id ShareId, input DeletedShare) (result RestoreOperationResponse, err error) {
	req, err := c.preparerForRestore(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Restore", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Restore", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRestore(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Restore", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRestore prepares the Restore request.
func (c FileSharesClient) preparerForRestore(ctx context.Context, id ShareId, input DeletedShare) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restore", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRestore handles the response to the Restore request. The method always
// closes the http.Response Body.
func (c FileSharesClient) responderForRestore(resp *http.Response) (result RestoreOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
