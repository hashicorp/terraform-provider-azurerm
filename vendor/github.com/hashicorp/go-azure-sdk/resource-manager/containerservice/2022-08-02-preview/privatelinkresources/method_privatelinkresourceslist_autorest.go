package privatelinkresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkResourcesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResourcesListResult
}

// PrivateLinkResourcesList ...
func (c PrivateLinkResourcesClient) PrivateLinkResourcesList(ctx context.Context, id ManagedClusterId) (result PrivateLinkResourcesListOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkResourcesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "PrivateLinkResourcesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "PrivateLinkResourcesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkResourcesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "PrivateLinkResourcesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkResourcesList prepares the PrivateLinkResourcesList request.
func (c PrivateLinkResourcesClient) preparerForPrivateLinkResourcesList(ctx context.Context, id ManagedClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateLinkResourcesList handles the response to the PrivateLinkResourcesList request. The method always
// closes the http.Response Body.
func (c PrivateLinkResourcesClient) responderForPrivateLinkResourcesList(resp *http.Response) (result PrivateLinkResourcesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
