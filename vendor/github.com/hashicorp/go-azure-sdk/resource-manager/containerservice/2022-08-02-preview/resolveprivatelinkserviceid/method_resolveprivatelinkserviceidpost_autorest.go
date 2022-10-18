package resolveprivatelinkserviceid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResolvePrivateLinkServiceIdPOSTOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResource
}

// ResolvePrivateLinkServiceIdPOST ...
func (c ResolvePrivateLinkServiceIdClient) ResolvePrivateLinkServiceIdPOST(ctx context.Context, id ManagedClusterId, input PrivateLinkResource) (result ResolvePrivateLinkServiceIdPOSTOperationResponse, err error) {
	req, err := c.preparerForResolvePrivateLinkServiceIdPOST(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "ResolvePrivateLinkServiceIdPOST", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "ResolvePrivateLinkServiceIdPOST", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForResolvePrivateLinkServiceIdPOST(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "ResolvePrivateLinkServiceIdPOST", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForResolvePrivateLinkServiceIdPOST prepares the ResolvePrivateLinkServiceIdPOST request.
func (c ResolvePrivateLinkServiceIdClient) preparerForResolvePrivateLinkServiceIdPOST(ctx context.Context, id ManagedClusterId, input PrivateLinkResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resolvePrivateLinkServiceId", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForResolvePrivateLinkServiceIdPOST handles the response to the ResolvePrivateLinkServiceIdPOST request. The method always
// closes the http.Response Body.
func (c ResolvePrivateLinkServiceIdClient) responderForResolvePrivateLinkServiceIdPOST(resp *http.Response) (result ResolvePrivateLinkServiceIdPOSTOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
