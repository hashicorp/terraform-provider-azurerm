package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOSOptionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *OSOptionProfile
}

type GetOSOptionsOperationOptions struct {
	ResourceType *string
}

func DefaultGetOSOptionsOperationOptions() GetOSOptionsOperationOptions {
	return GetOSOptionsOperationOptions{}
}

func (o GetOSOptionsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o GetOSOptionsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ResourceType != nil {
		out["resource-type"] = *o.ResourceType
	}

	return out
}

// GetOSOptions ...
func (c ManagedClustersClient) GetOSOptions(ctx context.Context, id LocationId, options GetOSOptionsOperationOptions) (result GetOSOptionsOperationResponse, err error) {
	req, err := c.preparerForGetOSOptions(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetOSOptions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetOSOptions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetOSOptions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "GetOSOptions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetOSOptions prepares the GetOSOptions request.
func (c ManagedClustersClient) preparerForGetOSOptions(ctx context.Context, id LocationId, options GetOSOptionsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/osOptions/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetOSOptions handles the response to the GetOSOptions request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForGetOSOptions(resp *http.Response) (result GetOSOptionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
