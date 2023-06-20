package managedidentities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedIdentitiesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// UserAssignedIdentitiesDelete ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesDelete(ctx context.Context, id commonids.UserAssignedIdentityId) (result UserAssignedIdentitiesDeleteOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesDelete prepares the UserAssignedIdentitiesDelete request.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesDelete(ctx context.Context, id commonids.UserAssignedIdentityId) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesDelete handles the response to the UserAssignedIdentitiesDelete request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForUserAssignedIdentitiesDelete(resp *http.Response) (result UserAssignedIdentitiesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
