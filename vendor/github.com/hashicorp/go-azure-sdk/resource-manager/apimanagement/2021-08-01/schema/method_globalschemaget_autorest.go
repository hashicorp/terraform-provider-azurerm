package schema

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalSchemaGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *GlobalSchemaContract
}

// GlobalSchemaGet ...
func (c SchemaClient) GlobalSchemaGet(ctx context.Context, id SchemaId) (result GlobalSchemaGetOperationResponse, err error) {
	req, err := c.preparerForGlobalSchemaGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGlobalSchemaGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGlobalSchemaGet prepares the GlobalSchemaGet request.
func (c SchemaClient) preparerForGlobalSchemaGet(ctx context.Context, id SchemaId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGlobalSchemaGet handles the response to the GlobalSchemaGet request. The method always
// closes the http.Response Body.
func (c SchemaClient) responderForGlobalSchemaGet(resp *http.Response) (result GlobalSchemaGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
