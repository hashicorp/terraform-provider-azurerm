package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoteRenderingAccountsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// RemoteRenderingAccountsDelete ...
func (c ResourceClient) RemoteRenderingAccountsDelete(ctx context.Context, id RemoteRenderingAccountId) (result RemoteRenderingAccountsDeleteOperationResponse, err error) {
	req, err := c.preparerForRemoteRenderingAccountsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemoteRenderingAccountsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemoteRenderingAccountsDelete prepares the RemoteRenderingAccountsDelete request.
func (c ResourceClient) preparerForRemoteRenderingAccountsDelete(ctx context.Context, id RemoteRenderingAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemoteRenderingAccountsDelete handles the response to the RemoteRenderingAccountsDelete request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForRemoteRenderingAccountsDelete(resp *http.Response) (result RemoteRenderingAccountsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
