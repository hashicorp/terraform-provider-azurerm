package availabilitysets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAvailableSizesOperationResponse struct {
	HttpResponse *http.Response
	Model        *VirtualMachineSizeListResult
}

// ListAvailableSizes ...
func (c AvailabilitySetsClient) ListAvailableSizes(ctx context.Context, id AvailabilitySetId) (result ListAvailableSizesOperationResponse, err error) {
	req, err := c.preparerForListAvailableSizes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitysets.AvailabilitySetsClient", "ListAvailableSizes", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitysets.AvailabilitySetsClient", "ListAvailableSizes", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListAvailableSizes(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitysets.AvailabilitySetsClient", "ListAvailableSizes", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListAvailableSizes prepares the ListAvailableSizes request.
func (c AvailabilitySetsClient) preparerForListAvailableSizes(ctx context.Context, id AvailabilitySetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/vmSizes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListAvailableSizes handles the response to the ListAvailableSizes request. The method always
// closes the http.Response Body.
func (c AvailabilitySetsClient) responderForListAvailableSizes(resp *http.Response) (result ListAvailableSizesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
