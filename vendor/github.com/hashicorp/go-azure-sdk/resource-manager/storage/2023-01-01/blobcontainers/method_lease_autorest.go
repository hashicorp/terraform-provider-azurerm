package blobcontainers

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

type LeaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *LeaseContainerResponse
}

// Lease ...
func (c BlobContainersClient) Lease(ctx context.Context, id commonids.StorageContainerId, input LeaseContainerRequest) (result LeaseOperationResponse, err error) {
	req, err := c.preparerForLease(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "Lease", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "Lease", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLease(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "Lease", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLease prepares the Lease request.
func (c BlobContainersClient) preparerForLease(ctx context.Context, id commonids.StorageContainerId, input LeaseContainerRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/lease", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLease handles the response to the Lease request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForLease(resp *http.Response) (result LeaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
