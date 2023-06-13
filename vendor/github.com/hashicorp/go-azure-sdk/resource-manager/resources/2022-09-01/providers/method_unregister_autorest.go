package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnregisterOperationResponse struct {
	HttpResponse *http.Response
	Model        *Provider
}

// Unregister ...
func (c ProvidersClient) Unregister(ctx context.Context, id SubscriptionProviderId) (result UnregisterOperationResponse, err error) {
	req, err := c.preparerForUnregister(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "Unregister", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "Unregister", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUnregister(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "Unregister", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUnregister prepares the Unregister request.
func (c ProvidersClient) preparerForUnregister(ctx context.Context, id SubscriptionProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/unregister", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUnregister handles the response to the Unregister request. The method always
// closes the http.Response Body.
func (c ProvidersClient) responderForUnregister(resp *http.Response) (result UnregisterOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
