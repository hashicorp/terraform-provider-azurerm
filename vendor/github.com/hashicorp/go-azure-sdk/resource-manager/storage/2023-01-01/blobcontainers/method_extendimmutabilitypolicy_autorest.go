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

type ExtendImmutabilityPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ImmutabilityPolicy
}

type ExtendImmutabilityPolicyOperationOptions struct {
	IfMatch *string
}

func DefaultExtendImmutabilityPolicyOperationOptions() ExtendImmutabilityPolicyOperationOptions {
	return ExtendImmutabilityPolicyOperationOptions{}
}

func (o ExtendImmutabilityPolicyOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o ExtendImmutabilityPolicyOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// ExtendImmutabilityPolicy ...
func (c BlobContainersClient) ExtendImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, input ImmutabilityPolicy, options ExtendImmutabilityPolicyOperationOptions) (result ExtendImmutabilityPolicyOperationResponse, err error) {
	req, err := c.preparerForExtendImmutabilityPolicy(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ExtendImmutabilityPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ExtendImmutabilityPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExtendImmutabilityPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ExtendImmutabilityPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExtendImmutabilityPolicy prepares the ExtendImmutabilityPolicy request.
func (c BlobContainersClient) preparerForExtendImmutabilityPolicy(ctx context.Context, id commonids.StorageContainerId, input ImmutabilityPolicy, options ExtendImmutabilityPolicyOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/immutabilityPolicies/default/extend", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForExtendImmutabilityPolicy handles the response to the ExtendImmutabilityPolicy request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForExtendImmutabilityPolicy(resp *http.Response) (result ExtendImmutabilityPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
