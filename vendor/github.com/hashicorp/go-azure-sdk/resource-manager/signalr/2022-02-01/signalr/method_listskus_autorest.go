package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSkusOperationResponse struct {
	HttpResponse *http.Response
	Model        *SkuList
}

// ListSkus ...
func (c SignalRClient) ListSkus(ctx context.Context, id SignalRId) (result ListSkusOperationResponse, err error) {
	req, err := c.preparerForListSkus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "ListSkus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "ListSkus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSkus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "ListSkus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSkus prepares the ListSkus request.
func (c SignalRClient) preparerForListSkus(ctx context.Context, id SignalRId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSkus handles the response to the ListSkus request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForListSkus(resp *http.Response) (result ListSkusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
