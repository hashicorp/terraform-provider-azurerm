package schema

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalSchemaDeleteOperationResponse struct {
	HttpResponse *http.Response
}

type GlobalSchemaDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultGlobalSchemaDeleteOperationOptions() GlobalSchemaDeleteOperationOptions {
	return GlobalSchemaDeleteOperationOptions{}
}

func (o GlobalSchemaDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o GlobalSchemaDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// GlobalSchemaDelete ...
func (c SchemaClient) GlobalSchemaDelete(ctx context.Context, id SchemaId, options GlobalSchemaDeleteOperationOptions) (result GlobalSchemaDeleteOperationResponse, err error) {
	req, err := c.preparerForGlobalSchemaDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGlobalSchemaDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGlobalSchemaDelete prepares the GlobalSchemaDelete request.
func (c SchemaClient) preparerForGlobalSchemaDelete(ctx context.Context, id SchemaId, options GlobalSchemaDeleteOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGlobalSchemaDelete handles the response to the GlobalSchemaDelete request. The method always
// closes the http.Response Body.
func (c SchemaClient) responderForGlobalSchemaDelete(resp *http.Response) (result GlobalSchemaDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
