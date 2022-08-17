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

type IntelligencePacksEnableOperationResponse struct {
	HttpResponse *http.Response
}

// IntelligencePacksEnable ...
func (c WorkspacesClient) IntelligencePacksEnable(ctx context.Context, id IntelligencePackId) (result IntelligencePacksEnableOperationResponse, err error) {
	req, err := c.preparerForIntelligencePacksEnable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksEnable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksEnable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIntelligencePacksEnable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "IntelligencePacksEnable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIntelligencePacksEnable prepares the IntelligencePacksEnable request.
func (c WorkspacesClient) preparerForIntelligencePacksEnable(ctx context.Context, id IntelligencePackId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForIntelligencePacksEnable handles the response to the IntelligencePacksEnable request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForIntelligencePacksEnable(resp *http.Response) (result IntelligencePacksEnableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
