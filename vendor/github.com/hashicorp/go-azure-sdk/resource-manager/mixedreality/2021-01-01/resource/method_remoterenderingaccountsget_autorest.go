package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoteRenderingAccountsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *RemoteRenderingAccount
}

// RemoteRenderingAccountsGet ...
func (c ResourceClient) RemoteRenderingAccountsGet(ctx context.Context, id RemoteRenderingAccountId) (result RemoteRenderingAccountsGetOperationResponse, err error) {
	req, err := c.preparerForRemoteRenderingAccountsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemoteRenderingAccountsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemoteRenderingAccountsGet prepares the RemoteRenderingAccountsGet request.
func (c ResourceClient) preparerForRemoteRenderingAccountsGet(ctx context.Context, id RemoteRenderingAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemoteRenderingAccountsGet handles the response to the RemoteRenderingAccountsGet request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForRemoteRenderingAccountsGet(resp *http.Response) (result RemoteRenderingAccountsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
