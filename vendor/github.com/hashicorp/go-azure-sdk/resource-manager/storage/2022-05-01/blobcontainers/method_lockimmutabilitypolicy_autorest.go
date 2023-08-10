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

type LockImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ImmutabilityPolicy
}

type LockImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultLockImmutabilityPolicyOperationOptions() LockImmutabilityPolicyOperationOptions {
	return LockImmutabilityPolicyOperationOptions{}
}

func (o LockImmutabilityPolicyOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o LockImmutabilityPolicyOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// LockImmutabilityPolicy ...
func (c BlobContainersClient) LockImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options LockImmutabilityPolicyOperationOptions) (result LockImmutabilityPolicyOperationResponse, err error) {
	req, err := c.preparerForLockImmutabilityPolicy(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "LockImmutabilityPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "LockImmutabilityPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLockImmutabilityPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "LockImmutabilityPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLockImmutabilityPolicy prepares the LockImmutabilityPolicy request.
func (c BlobContainersClient) preparerForLockImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options LockImmutabilityPolicyOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/immutabilityPolicies/default/lock", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLockImmutabilityPolicy handles the response to the LockImmutabilityPolicy request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForLockImmutabilityPolicy(resp *http.Response) (result LockImmutabilityPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
