package workbooksapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// WorkbooksDelete ...
func (c WorkbooksAPIsClient) WorkbooksDelete(ctx context.Context, id WorkbookId) (result WorkbooksDeleteOperationResponse, err error) {
	req, err := c.preparerForWorkbooksDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbooksDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbooksDelete prepares the WorkbooksDelete request.
func (c WorkbooksAPIsClient) preparerForWorkbooksDelete(ctx context.Context, id WorkbookId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbooksDelete handles the response to the WorkbooksDelete request. The method always
// closes the http.Response Body.
func (c WorkbooksAPIsClient) responderForWorkbooksDelete(resp *http.Response) (result WorkbooksDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
