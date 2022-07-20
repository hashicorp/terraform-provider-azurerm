package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoteRenderingAccountsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *RemoteRenderingAccount
}

// RemoteRenderingAccountsCreate ...
func (c ResourceClient) RemoteRenderingAccountsCreate(ctx context.Context, id RemoteRenderingAccountId, input RemoteRenderingAccount) (result RemoteRenderingAccountsCreateOperationResponse, err error) {
	req, err := c.preparerForRemoteRenderingAccountsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemoteRenderingAccountsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemoteRenderingAccountsCreate prepares the RemoteRenderingAccountsCreate request.
func (c ResourceClient) preparerForRemoteRenderingAccountsCreate(ctx context.Context, id RemoteRenderingAccountId, input RemoteRenderingAccount) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemoteRenderingAccountsCreate handles the response to the RemoteRenderingAccountsCreate request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForRemoteRenderingAccountsCreate(resp *http.Response) (result RemoteRenderingAccountsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
