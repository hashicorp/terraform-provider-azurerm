package managedidentity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemAssignedIdentitiesGetByScopeOperationResponse struct {
	HttpResponse *http.Response
	Model        *SystemAssignedIdentity
}

// SystemAssignedIdentitiesGetByScope ...
func (c ManagedIdentityClient) SystemAssignedIdentitiesGetByScope(ctx context.Context, id commonids.ScopeId) (result SystemAssignedIdentitiesGetByScopeOperationResponse, err error) {
	req, err := c.preparerForSystemAssignedIdentitiesGetByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "SystemAssignedIdentitiesGetByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "SystemAssignedIdentitiesGetByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSystemAssignedIdentitiesGetByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "SystemAssignedIdentitiesGetByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSystemAssignedIdentitiesGetByScope prepares the SystemAssignedIdentitiesGetByScope request.
func (c ManagedIdentityClient) preparerForSystemAssignedIdentitiesGetByScope(ctx context.Context, id commonids.ScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ManagedIdentity/identities/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSystemAssignedIdentitiesGetByScope handles the response to the SystemAssignedIdentitiesGetByScope request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForSystemAssignedIdentitiesGetByScope(resp *http.Response) (result SystemAssignedIdentitiesGetByScopeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
