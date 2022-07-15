package workspaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListNotebookAccessTokenOperationResponse struct {
	HttpResponse *http.Response
	Model        *NotebookAccessTokenResult
}

// ListNotebookAccessToken ...
func (c WorkspacesClient) ListNotebookAccessToken(ctx context.Context, id WorkspaceId) (result ListNotebookAccessTokenOperationResponse, err error) {
	req, err := c.preparerForListNotebookAccessToken(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ListNotebookAccessToken", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ListNotebookAccessToken", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListNotebookAccessToken(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ListNotebookAccessToken", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListNotebookAccessToken prepares the ListNotebookAccessToken request.
func (c WorkspacesClient) preparerForListNotebookAccessToken(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listNotebookAccessToken", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListNotebookAccessToken handles the response to the ListNotebookAccessToken request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForListNotebookAccessToken(resp *http.Response) (result ListNotebookAccessTokenOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
