package managedidentities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FederatedIdentityCredentialsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// FederatedIdentityCredentialsDelete ...
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsDelete(ctx context.Context, id FederatedIdentityCredentialId) (result FederatedIdentityCredentialsDeleteOperationResponse, err error) {
	req, err := c.preparerForFederatedIdentityCredentialsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFederatedIdentityCredentialsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFederatedIdentityCredentialsDelete prepares the FederatedIdentityCredentialsDelete request.
func (c ManagedIdentitiesClient) preparerForFederatedIdentityCredentialsDelete(ctx context.Context, id FederatedIdentityCredentialId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForFederatedIdentityCredentialsDelete handles the response to the FederatedIdentityCredentialsDelete request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForFederatedIdentityCredentialsDelete(resp *http.Response) (result FederatedIdentityCredentialsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
