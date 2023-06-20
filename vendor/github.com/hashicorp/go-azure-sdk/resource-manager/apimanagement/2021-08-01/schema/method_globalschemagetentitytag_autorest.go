package schema

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalSchemaGetEntityTagOperationResponse struct {
	HttpResponse *http.Response
}

// GlobalSchemaGetEntityTag ...
func (c SchemaClient) GlobalSchemaGetEntityTag(ctx context.Context, id SchemaId) (result GlobalSchemaGetEntityTagOperationResponse, err error) {
	req, err := c.preparerForGlobalSchemaGetEntityTag(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGetEntityTag", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGetEntityTag", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGlobalSchemaGetEntityTag(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGetEntityTag", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGlobalSchemaGetEntityTag prepares the GlobalSchemaGetEntityTag request.
func (c SchemaClient) preparerForGlobalSchemaGetEntityTag(ctx context.Context, id SchemaId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsHead(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGlobalSchemaGetEntityTag handles the response to the GlobalSchemaGetEntityTag request. The method always
// closes the http.Response Body.
func (c SchemaClient) responderForGlobalSchemaGetEntityTag(resp *http.Response) (result GlobalSchemaGetEntityTagOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
