package resourceguardproxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnlockDeleteOperationResponse struct {
	HttpResponse *http.Response
	Model        *UnlockDeleteResponse
}

// UnlockDelete ...
func (c ResourceGuardProxyClient) UnlockDelete(ctx context.Context, id BackupResourceGuardProxyId, input UnlockDeleteRequest) (result UnlockDeleteOperationResponse, err error) {
	req, err := c.preparerForUnlockDelete(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxy.ResourceGuardProxyClient", "UnlockDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxy.ResourceGuardProxyClient", "UnlockDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUnlockDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxy.ResourceGuardProxyClient", "UnlockDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUnlockDelete prepares the UnlockDelete request.
func (c ResourceGuardProxyClient) preparerForUnlockDelete(ctx context.Context, id BackupResourceGuardProxyId, input UnlockDeleteRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/unlockDelete", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUnlockDelete handles the response to the UnlockDelete request. The method always
// closes the http.Response Body.
func (c ResourceGuardProxyClient) responderForUnlockDelete(resp *http.Response) (result UnlockDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
