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

type IntelligencePacksDisableOperationResponse struct {
	HttpResponse *http.Response
}

// IntelligencePacksDisable ...
func (c WorkspacesClient) IntelligencePacksDisable(ctx context.Context, id IntelligencePackId) (result IntelligencePacksDisableOperationResponse, err error) {
	req, err := c.preparerForIntelligencePacksDisable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksDisable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksDisable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIntelligencePacksDisable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksDisable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIntelligencePacksDisable prepares the IntelligencePacksDisable request.
func (c WorkspacesClient) preparerForIntelligencePacksDisable(ctx context.Context, id IntelligencePackId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIntelligencePacksDisable handles the response to the IntelligencePacksDisable request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForIntelligencePacksDisable(resp *http.Response) (result IntelligencePacksDisableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
