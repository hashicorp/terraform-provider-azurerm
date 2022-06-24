package operationstatus

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGroupContextGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *OperationResource
}

// ResourceGroupContextGet ...
func (c OperationStatusClient) ResourceGroupContextGet(ctx context.Context, id ProviderOperationStatuId) (result ResourceGroupContextGetOperationResponse, err error) {
	req, err := c.preparerForResourceGroupContextGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "ResourceGroupContextGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "ResourceGroupContextGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForResourceGroupContextGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationstatus.OperationStatusClient", "ResourceGroupContextGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForResourceGroupContextGet prepares the ResourceGroupContextGet request.
func (c OperationStatusClient) preparerForResourceGroupContextGet(ctx context.Context, id ProviderOperationStatuId) (*http.Request, error) {
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

// responderForResourceGroupContextGet handles the response to the ResourceGroupContextGet request. The method always
// closes the http.Response Body.
func (c OperationStatusClient) responderForResourceGroupContextGet(resp *http.Response) (result ResourceGroupContextGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
