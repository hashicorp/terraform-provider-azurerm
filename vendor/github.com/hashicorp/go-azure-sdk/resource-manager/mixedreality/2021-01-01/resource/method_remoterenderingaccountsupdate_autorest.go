package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoteRenderingAccountsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *RemoteRenderingAccount
}

// RemoteRenderingAccountsUpdate ...
func (c ResourceClient) RemoteRenderingAccountsUpdate(ctx context.Context, id RemoteRenderingAccountId, input RemoteRenderingAccount) (result RemoteRenderingAccountsUpdateOperationResponse, err error) {
	req, err := c.preparerForRemoteRenderingAccountsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemoteRenderingAccountsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "RemoteRenderingAccountsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemoteRenderingAccountsUpdate prepares the RemoteRenderingAccountsUpdate request.
func (c ResourceClient) preparerForRemoteRenderingAccountsUpdate(ctx context.Context, id RemoteRenderingAccountId, input RemoteRenderingAccount) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemoteRenderingAccountsUpdate handles the response to the RemoteRenderingAccountsUpdate request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForRemoteRenderingAccountsUpdate(resp *http.Response) (result RemoteRenderingAccountsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
