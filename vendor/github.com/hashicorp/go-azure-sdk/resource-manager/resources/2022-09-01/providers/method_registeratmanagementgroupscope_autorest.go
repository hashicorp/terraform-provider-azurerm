package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisterAtManagementGroupScopeOperationResponse struct {
	HttpResponse *http.Response
}

// RegisterAtManagementGroupScope ...
func (c ProvidersClient) RegisterAtManagementGroupScope(ctx context.Context, id Providers2Id) (result RegisterAtManagementGroupScopeOperationResponse, err error) {
	req, err := c.preparerForRegisterAtManagementGroupScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "RegisterAtManagementGroupScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "RegisterAtManagementGroupScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegisterAtManagementGroupScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "RegisterAtManagementGroupScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegisterAtManagementGroupScope prepares the RegisterAtManagementGroupScope request.
func (c ProvidersClient) preparerForRegisterAtManagementGroupScope(ctx context.Context, id Providers2Id) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/register", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegisterAtManagementGroupScope handles the response to the RegisterAtManagementGroupScope request. The method always
// closes the http.Response Body.
func (c ProvidersClient) responderForRegisterAtManagementGroupScope(resp *http.Response) (result RegisterAtManagementGroupScopeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
