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

type UserAssignedIdentitiesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Identity
}

// UserAssignedIdentitiesGet ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesGet(ctx context.Context, id commonids.UserAssignedIdentityId) (result UserAssignedIdentitiesGetOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesGet prepares the UserAssignedIdentitiesGet request.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesGet(ctx context.Context, id commonids.UserAssignedIdentityId) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesGet handles the response to the UserAssignedIdentitiesGet request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForUserAssignedIdentitiesGet(resp *http.Response) (result UserAssignedIdentitiesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
