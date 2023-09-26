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

type DeleteImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ImmutabilityPolicy
}

type DeleteImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultDeleteImmutabilityPolicyOperationOptions() DeleteImmutabilityPolicyOperationOptions {
	return DeleteImmutabilityPolicyOperationOptions{}
}

func (o DeleteImmutabilityPolicyOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o DeleteImmutabilityPolicyOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// DeleteImmutabilityPolicy ...
func (c BlobContainersClient) DeleteImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options DeleteImmutabilityPolicyOperationOptions) (result DeleteImmutabilityPolicyOperationResponse, err error) {
	req, err := c.preparerForDeleteImmutabilityPolicy(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "DeleteImmutabilityPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "DeleteImmutabilityPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteImmutabilityPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "DeleteImmutabilityPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteImmutabilityPolicy prepares the DeleteImmutabilityPolicy request.
func (c BlobContainersClient) preparerForDeleteImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, options DeleteImmutabilityPolicyOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/immutabilityPolicies/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeleteImmutabilityPolicy handles the response to the DeleteImmutabilityPolicy request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForDeleteImmutabilityPolicy(resp *http.Response) (result DeleteImmutabilityPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
