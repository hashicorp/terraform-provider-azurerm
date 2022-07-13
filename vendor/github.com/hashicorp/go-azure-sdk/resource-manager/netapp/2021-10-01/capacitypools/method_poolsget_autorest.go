package capacitypools

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *CapacityPool
}

// PoolsGet ...
func (c CapacityPoolsClient) PoolsGet(ctx context.Context, id CapacityPoolId) (result PoolsGetOperationResponse, err error) {
	req, err := c.preparerForPoolsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPoolsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPoolsGet prepares the PoolsGet request.
func (c CapacityPoolsClient) preparerForPoolsGet(ctx context.Context, id CapacityPoolId) (*http.Request, error) {
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

// responderForPoolsGet handles the response to the PoolsGet request. The method always
// closes the http.Response Body.
func (c CapacityPoolsClient) responderForPoolsGet(resp *http.Response) (result PoolsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
