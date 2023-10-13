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

type GetImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ImmutabilityPolicy
}

type GetImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultGetImmutabilityPolicyOperationOptions() GetImmutabilityPolicyOperationOptions {
	return GetImmutabilityPolicyOperationOptions{}
}

func (o GetImmutabilityPolicyOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o GetImmutabilityPolicyOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// GetImmutabilityPolicy ...
func (c BlobContainersClient) GetImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options GetImmutabilityPolicyOperationOptions) (result GetImmutabilityPolicyOperationResponse, err error) {
	req, err := c.preparerForGetImmutabilityPolicy(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "GetImmutabilityPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "GetImmutabilityPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetImmutabilityPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "GetImmutabilityPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetImmutabilityPolicy prepares the GetImmutabilityPolicy request.
func (c BlobContainersClient) preparerForGetImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options GetImmutabilityPolicyOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/immutabilityPolicies/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetImmutabilityPolicy handles the response to the GetImmutabilityPolicy request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForGetImmutabilityPolicy(resp *http.Response) (result GetImmutabilityPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
