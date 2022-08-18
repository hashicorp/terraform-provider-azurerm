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

type IntelligencePacksListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]IntelligencePack
}

// IntelligencePacksList ...
func (c WorkspacesClient) IntelligencePacksList(ctx context.Context, id WorkspaceId) (result IntelligencePacksListOperationResponse, err error) {
	req, err := c.preparerForIntelligencePacksList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIntelligencePacksList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIntelligencePacksList prepares the IntelligencePacksList request.
func (c WorkspacesClient) preparerForIntelligencePacksList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/intelligencePacks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIntelligencePacksList handles the response to the IntelligencePacksList request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForIntelligencePacksList(resp *http.Response) (result IntelligencePacksListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
