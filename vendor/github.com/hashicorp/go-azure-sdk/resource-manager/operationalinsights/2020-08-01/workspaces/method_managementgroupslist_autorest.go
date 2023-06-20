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

type ManagementGroupsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkspaceListManagementGroupsResult
}

// ManagementGroupsList ...
func (c WorkspacesClient) ManagementGroupsList(ctx context.Context, id WorkspaceId) (result ManagementGroupsListOperationResponse, err error) {
	req, err := c.preparerForManagementGroupsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ManagementGroupsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ManagementGroupsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagementGroupsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "ManagementGroupsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagementGroupsList prepares the ManagementGroupsList request.
func (c WorkspacesClient) preparerForManagementGroupsList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/managementGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForManagementGroupsList handles the response to the ManagementGroupsList request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForManagementGroupsList(resp *http.Response) (result ManagementGroupsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
