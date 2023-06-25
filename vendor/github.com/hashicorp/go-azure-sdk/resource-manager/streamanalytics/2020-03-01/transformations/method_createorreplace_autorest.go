package transformations

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrReplaceOperationResponse struct {
	HttpResponse *http.Response
	Model        *Transformation
}

type CreateOrReplaceOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOrReplaceOperationOptions() CreateOrReplaceOperationOptions {
	return CreateOrReplaceOperationOptions{}
}

func (o CreateOrReplaceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	if o.IfNoneMatch != nil {
		out["If-None-Match"] = *o.IfNoneMatch
	}

	return out
}

func (o CreateOrReplaceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// CreateOrReplace ...
func (c TransformationsClient) CreateOrReplace(ctx context.Context, id TransformationId, input Transformation, options CreateOrReplaceOperationOptions) (result CreateOrReplaceOperationResponse, err error) {
	req, err := c.preparerForCreateOrReplace(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "transformations.TransformationsClient", "CreateOrReplace", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "transformations.TransformationsClient", "CreateOrReplace", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrReplace(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "transformations.TransformationsClient", "CreateOrReplace", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrReplace prepares the CreateOrReplace request.
func (c TransformationsClient) preparerForCreateOrReplace(ctx context.Context, id TransformationId, input Transformation, options CreateOrReplaceOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreateOrReplace handles the response to the CreateOrReplace request. The method always
// closes the http.Response Body.
func (c TransformationsClient) responderForCreateOrReplace(resp *http.Response) (result CreateOrReplaceOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
