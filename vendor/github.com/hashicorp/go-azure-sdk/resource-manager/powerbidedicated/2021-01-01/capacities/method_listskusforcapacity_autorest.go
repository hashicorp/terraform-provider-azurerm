package capacities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSkusForCapacityOperationResponse struct {
	HttpResponse *http.Response
	Model        *SkuEnumerationForExistingResourceResult
}

// ListSkusForCapacity ...
func (c CapacitiesClient) ListSkusForCapacity(ctx context.Context, id CapacitiesId) (result ListSkusForCapacityOperationResponse, err error) {
	req, err := c.preparerForListSkusForCapacity(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacities.CapacitiesClient", "ListSkusForCapacity", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacities.CapacitiesClient", "ListSkusForCapacity", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSkusForCapacity(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacities.CapacitiesClient", "ListSkusForCapacity", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSkusForCapacity prepares the ListSkusForCapacity request.
func (c CapacitiesClient) preparerForListSkusForCapacity(ctx context.Context, id CapacitiesId) (*http.Request, error) {
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

// responderForListSkusForCapacity handles the response to the ListSkusForCapacity request. The method always
// closes the http.Response Body.
func (c CapacitiesClient) responderForListSkusForCapacity(resp *http.Response) (result ListSkusForCapacityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
