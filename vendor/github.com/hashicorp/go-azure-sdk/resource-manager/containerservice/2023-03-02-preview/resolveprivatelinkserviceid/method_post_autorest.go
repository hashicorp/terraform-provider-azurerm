package resolveprivatelinkserviceid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type POSTOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResource
}

// POST ...
func (c ResolvePrivateLinkServiceIdClient) POST(ctx context.Context, id commonids.KubernetesClusterId, input PrivateLinkResource) (result POSTOperationResponse, err error) {
	req, err := c.preparerForPOST(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "POST", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "POST", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPOST(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient", "POST", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPOST prepares the POST request.
func (c ResolvePrivateLinkServiceIdClient) preparerForPOST(ctx context.Context, id commonids.KubernetesClusterId, input PrivateLinkResource) (*http.Request, error) {
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

// responderForPOST handles the response to the POST request. The method always
// closes the http.Response Body.
func (c ResolvePrivateLinkServiceIdClient) responderForPOST(resp *http.Response) (result POSTOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
