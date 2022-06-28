package adminkeys

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegenerateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AdminKeyResult
}

type RegenerateOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultRegenerateOperationOptions() RegenerateOperationOptions {
	return RegenerateOperationOptions{}
}

func (o RegenerateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsClientRequestId != nil {
		out["x-ms-client-request-id"] = *o.XMsClientRequestId
	}

	return out
}

func (o RegenerateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Regenerate ...
func (c AdminKeysClient) Regenerate(ctx context.Context, id KeyKindId, options RegenerateOperationOptions) (result RegenerateOperationResponse, err error) {
	req, err := c.preparerForRegenerate(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "adminkeys.AdminKeysClient", "Regenerate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "adminkeys.AdminKeysClient", "Regenerate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegenerate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "adminkeys.AdminKeysClient", "Regenerate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegenerate prepares the Regenerate request.
func (c AdminKeysClient) preparerForRegenerate(ctx context.Context, id KeyKindId, options RegenerateOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegenerate handles the response to the Regenerate request. The method always
// closes the http.Response Body.
func (c AdminKeysClient) responderForRegenerate(resp *http.Response) (result RegenerateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
