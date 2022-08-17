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

type WorkspacePurgePurgeOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkspacePurgeResponse
}

// WorkspacePurgePurge ...
func (c WorkspacesClient) WorkspacePurgePurge(ctx context.Context, id WorkspaceId, input WorkspacePurgeBody) (result WorkspacePurgePurgeOperationResponse, err error) {
	req, err := c.preparerForWorkspacePurgePurge(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgePurge", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgePurge", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkspacePurgePurge(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspacePurgePurge", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkspacePurgePurge prepares the WorkspacePurgePurge request.
func (c WorkspacesClient) preparerForWorkspacePurgePurge(ctx context.Context, id WorkspaceId, input WorkspacePurgeBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/purge", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkspacePurgePurge handles the response to the WorkspacePurgePurge request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForWorkspacePurgePurge(resp *http.Response) (result WorkspacePurgePurgeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
