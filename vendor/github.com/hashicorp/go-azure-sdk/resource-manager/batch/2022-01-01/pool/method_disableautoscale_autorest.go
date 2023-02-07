package pool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisableAutoScaleOperationResponse struct {
	HttpResponse *http.Response
	Model        *Pool
}

// DisableAutoScale ...
func (c PoolClient) DisableAutoScale(ctx context.Context, id PoolId) (result DisableAutoScaleOperationResponse, err error) {
	req, err := c.preparerForDisableAutoScale(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "DisableAutoScale", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "DisableAutoScale", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisableAutoScale(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "DisableAutoScale", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisableAutoScale prepares the DisableAutoScale request.
func (c PoolClient) preparerForDisableAutoScale(ctx context.Context, id PoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disableAutoScale", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDisableAutoScale handles the response to the DisableAutoScale request. The method always
// closes the http.Response Body.
func (c PoolClient) responderForDisableAutoScale(resp *http.Response) (result DisableAutoScaleOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
