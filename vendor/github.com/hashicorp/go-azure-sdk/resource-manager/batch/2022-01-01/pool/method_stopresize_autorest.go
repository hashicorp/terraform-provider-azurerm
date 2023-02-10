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

type StopResizeOperationResponse struct {
	HttpResponse *http.Response
	Model        *Pool
}

// StopResize ...
func (c PoolClient) StopResize(ctx context.Context, id PoolId) (result StopResizeOperationResponse, err error) {
	req, err := c.preparerForStopResize(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "StopResize", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "StopResize", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStopResize(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "pool.PoolClient", "StopResize", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStopResize prepares the StopResize request.
func (c PoolClient) preparerForStopResize(ctx context.Context, id PoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stopResize", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStopResize handles the response to the StopResize request. The method always
// closes the http.Response Body.
func (c PoolClient) responderForStopResize(resp *http.Response) (result StopResizeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
