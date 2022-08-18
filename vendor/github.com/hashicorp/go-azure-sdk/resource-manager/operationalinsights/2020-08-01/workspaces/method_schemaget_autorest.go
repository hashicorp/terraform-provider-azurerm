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

type SchemaGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SearchGetSchemaResponse
}

// SchemaGet ...
func (c WorkspacesClient) SchemaGet(ctx context.Context, id WorkspaceId) (result SchemaGetOperationResponse, err error) {
	req, err := c.preparerForSchemaGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SchemaGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SchemaGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSchemaGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "SchemaGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSchemaGet prepares the SchemaGet request.
func (c WorkspacesClient) preparerForSchemaGet(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/schema", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSchemaGet handles the response to the SchemaGet request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForSchemaGet(resp *http.Response) (result SchemaGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
