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

type CreateOrUpdateImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ImmutabilityPolicy
}

type CreateOrUpdateImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultCreateOrUpdateImmutabilityPolicyOperationOptions() CreateOrUpdateImmutabilityPolicyOperationOptions {
	return CreateOrUpdateImmutabilityPolicyOperationOptions{}
}

func (o CreateOrUpdateImmutabilityPolicyOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o CreateOrUpdateImmutabilityPolicyOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// CreateOrUpdateImmutabilityPolicy ...
func (c BlobContainersClient) CreateOrUpdateImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, input ImmutabilityPolicy, options CreateOrUpdateImmutabilityPolicyOperationOptions) (result CreateOrUpdateImmutabilityPolicyOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateImmutabilityPolicy(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "CreateOrUpdateImmutabilityPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "CreateOrUpdateImmutabilityPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateImmutabilityPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "CreateOrUpdateImmutabilityPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateImmutabilityPolicy prepares the CreateOrUpdateImmutabilityPolicy request.
func (c BlobContainersClient) preparerForCreateOrUpdateImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, input ImmutabilityPolicy, options CreateOrUpdateImmutabilityPolicyOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/immutabilityPolicies/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreateOrUpdateImmutabilityPolicy handles the response to the CreateOrUpdateImmutabilityPolicy request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForCreateOrUpdateImmutabilityPolicy(resp *http.Response) (result CreateOrUpdateImmutabilityPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
