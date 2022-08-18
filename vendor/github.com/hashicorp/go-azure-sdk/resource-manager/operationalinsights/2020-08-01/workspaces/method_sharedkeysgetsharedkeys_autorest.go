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

type SharedKeysGetSharedKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedKeys
}

// SharedKeysGetSharedKeys ...
func (c WorkspacesClient) SharedKeysGetSharedKeys(ctx context.Context, id WorkspaceId) (result SharedKeysGetSharedKeysOperationResponse, err error) {
	req, err := c.preparerForSharedKeysGetSharedKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysGetSharedKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysGetSharedKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSharedKeysGetSharedKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysGetSharedKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSharedKeysGetSharedKeys prepares the SharedKeysGetSharedKeys request.
func (c WorkspacesClient) preparerForSharedKeysGetSharedKeys(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/sharedKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSharedKeysGetSharedKeys handles the response to the SharedKeysGetSharedKeys request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForSharedKeysGetSharedKeys(resp *http.Response) (result SharedKeysGetSharedKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
