package iotdpsresource

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

type ListPrivateLinkResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResources
}

// ListPrivateLinkResources ...
func (c IotDpsResourceClient) ListPrivateLinkResources(ctx context.Context, id commonids.ProvisioningServiceId) (result ListPrivateLinkResourcesOperationResponse, err error) {
	req, err := c.preparerForListPrivateLinkResources(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateLinkResources", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateLinkResources", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListPrivateLinkResources(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListPrivateLinkResources", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListPrivateLinkResources prepares the ListPrivateLinkResources request.
func (c IotDpsResourceClient) preparerForListPrivateLinkResources(ctx context.Context, id commonids.ProvisioningServiceId) (*http.Request, error) {
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

// responderForListPrivateLinkResources handles the response to the ListPrivateLinkResources request. The method always
// closes the http.Response Body.
func (c IotDpsResourceClient) responderForListPrivateLinkResources(resp *http.Response) (result ListPrivateLinkResourcesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
