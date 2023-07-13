package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderResourceTypesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ProviderResourceTypeListResult
}

type ProviderResourceTypesListOperationOptions struct {
	Expand *string
}

func DefaultProviderResourceTypesListOperationOptions() ProviderResourceTypesListOperationOptions {
	return ProviderResourceTypesListOperationOptions{}
}

func (o ProviderResourceTypesListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ProviderResourceTypesListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// ProviderResourceTypesList ...
func (c ProvidersClient) ProviderResourceTypesList(ctx context.Context, id SubscriptionProviderId, options ProviderResourceTypesListOperationOptions) (result ProviderResourceTypesListOperationResponse, err error) {
	req, err := c.preparerForProviderResourceTypesList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderResourceTypesList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderResourceTypesList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProviderResourceTypesList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderResourceTypesList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProviderResourceTypesList prepares the ProviderResourceTypesList request.
func (c ProvidersClient) preparerForProviderResourceTypesList(ctx context.Context, id SubscriptionProviderId, options ProviderResourceTypesListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/resourceTypes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForProviderResourceTypesList handles the response to the ProviderResourceTypesList request. The method always
// closes the http.Response Body.
func (c ProvidersClient) responderForProviderResourceTypesList(resp *http.Response) (result ProviderResourceTypesListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
