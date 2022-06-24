package workspaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePurgeGetPurgeStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkspacePurgeStatusResponse
}

// WorkspacePurgeGetPurgeStatus ...
func (c WorkspacesClient) WorkspacePurgeGetPurgeStatus(ctx context.Context, id OperationId) (result WorkspacePurgeGetPurgeStatusOperationResponse, err error) {
	req, err := c.preparerForWorkspacePurgeGetPurgeStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgeGetPurgeStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgeGetPurgeStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkspacePurgeGetPurgeStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgeGetPurgeStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkspacePurgeGetPurgeStatus prepares the WorkspacePurgeGetPurgeStatus request.
func (c WorkspacesClient) preparerForWorkspacePurgeGetPurgeStatus(ctx context.Context, id OperationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkspacePurgeGetPurgeStatus handles the response to the WorkspacePurgeGetPurgeStatus request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForWorkspacePurgeGetPurgeStatus(resp *http.Response) (result WorkspacePurgeGetPurgeStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
