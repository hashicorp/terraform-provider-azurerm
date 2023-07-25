package managedidentities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FederatedIdentityCredentialsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *FederatedIdentityCredential
}

// FederatedIdentityCredentialsCreateOrUpdate ...
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsCreateOrUpdate(ctx context.Context, id FederatedIdentityCredentialId, input FederatedIdentityCredential) (result FederatedIdentityCredentialsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForFederatedIdentityCredentialsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFederatedIdentityCredentialsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFederatedIdentityCredentialsCreateOrUpdate prepares the FederatedIdentityCredentialsCreateOrUpdate request.
func (c ManagedIdentitiesClient) preparerForFederatedIdentityCredentialsCreateOrUpdate(ctx context.Context, id FederatedIdentityCredentialId, input FederatedIdentityCredential) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForFederatedIdentityCredentialsCreateOrUpdate handles the response to the FederatedIdentityCredentialsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForFederatedIdentityCredentialsCreateOrUpdate(resp *http.Response) (result FederatedIdentityCredentialsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
