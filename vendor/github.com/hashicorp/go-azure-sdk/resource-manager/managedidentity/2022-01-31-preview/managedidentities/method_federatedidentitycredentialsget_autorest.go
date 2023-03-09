package managedidentities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FederatedIdentityCredentialsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *FederatedIdentityCredential
}

// FederatedIdentityCredentialsGet ...
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsGet(ctx context.Context, id FederatedIdentityCredentialId) (result FederatedIdentityCredentialsGetOperationResponse, err error) {
	req, err := c.preparerForFederatedIdentityCredentialsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFederatedIdentityCredentialsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFederatedIdentityCredentialsGet prepares the FederatedIdentityCredentialsGet request.
func (c ManagedIdentitiesClient) preparerForFederatedIdentityCredentialsGet(ctx context.Context, id FederatedIdentityCredentialId) (*http.Request, error) {
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

// responderForFederatedIdentityCredentialsGet handles the response to the FederatedIdentityCredentialsGet request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForFederatedIdentityCredentialsGet(resp *http.Response) (result FederatedIdentityCredentialsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
