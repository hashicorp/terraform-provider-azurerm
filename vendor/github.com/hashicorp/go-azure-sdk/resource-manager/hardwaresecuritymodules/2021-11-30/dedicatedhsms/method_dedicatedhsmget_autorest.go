package dedicatedhsms

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHsmGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DedicatedHsm
}

// DedicatedHsmGet ...
func (c DedicatedHsmsClient) DedicatedHsmGet(ctx context.Context, id DedicatedHSMId) (result DedicatedHsmGetOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDedicatedHsmGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDedicatedHsmGet prepares the DedicatedHsmGet request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmGet(ctx context.Context, id DedicatedHSMId) (*http.Request, error) {
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

// responderForDedicatedHsmGet handles the response to the DedicatedHsmGet request. The method always
// closes the http.Response Body.
func (c DedicatedHsmsClient) responderForDedicatedHsmGet(resp *http.Response) (result DedicatedHsmGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
