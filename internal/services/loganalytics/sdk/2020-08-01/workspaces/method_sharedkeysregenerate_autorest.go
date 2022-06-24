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

type SharedKeysRegenerateOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedKeys
}

// SharedKeysRegenerate ...
func (c WorkspacesClient) SharedKeysRegenerate(ctx context.Context, id WorkspaceId) (result SharedKeysRegenerateOperationResponse, err error) {
	req, err := c.preparerForSharedKeysRegenerate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysRegenerate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysRegenerate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSharedKeysRegenerate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SharedKeysRegenerate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSharedKeysRegenerate prepares the SharedKeysRegenerate request.
func (c WorkspacesClient) preparerForSharedKeysRegenerate(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateSharedKey", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSharedKeysRegenerate handles the response to the SharedKeysRegenerate request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForSharedKeysRegenerate(resp *http.Response) (result SharedKeysRegenerateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
