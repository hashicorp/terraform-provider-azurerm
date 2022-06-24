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

type UsagesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkspaceListUsagesResult
}

// UsagesList ...
func (c WorkspacesClient) UsagesList(ctx context.Context, id WorkspaceId) (result UsagesListOperationResponse, err error) {
	req, err := c.preparerForUsagesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "UsagesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "UsagesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUsagesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "UsagesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUsagesList prepares the UsagesList request.
func (c WorkspacesClient) preparerForUsagesList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/usages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUsagesList handles the response to the UsagesList request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForUsagesList(resp *http.Response) (result UsagesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
