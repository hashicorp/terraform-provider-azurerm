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

type ProviderPermissionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ProviderPermissionListResult
}

// ProviderPermissions ...
func (c ProvidersClient) ProviderPermissions(ctx context.Context, id SubscriptionProviderId) (result ProviderPermissionsOperationResponse, err error) {
	req, err := c.preparerForProviderPermissions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderPermissions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderPermissions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProviderPermissions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "providers.ProvidersClient", "ProviderPermissions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProviderPermissions prepares the ProviderPermissions request.
func (c ProvidersClient) preparerForProviderPermissions(ctx context.Context, id SubscriptionProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providerPermissions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForProviderPermissions handles the response to the ProviderPermissions request. The method always
// closes the http.Response Body.
func (c ProvidersClient) responderForProviderPermissions(resp *http.Response) (result ProviderPermissionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
