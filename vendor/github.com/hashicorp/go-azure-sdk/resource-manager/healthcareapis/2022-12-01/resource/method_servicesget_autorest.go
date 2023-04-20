package resource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ServicesDescription
}

// ServicesGet ...
func (c ResourceClient) ServicesGet(ctx context.Context, id ServiceId) (result ServicesGetOperationResponse, err error) {
	req, err := c.preparerForServicesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "ServicesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "ServicesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServicesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "ServicesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServicesGet prepares the ServicesGet request.
func (c ResourceClient) preparerForServicesGet(ctx context.Context, id ServiceId) (*http.Request, error) {
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

// responderForServicesGet handles the response to the ServicesGet request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForServicesGet(resp *http.Response) (result ServicesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
