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

type UserAssignedIdentitiesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Identity
}

// UserAssignedIdentitiesUpdate ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesUpdate(ctx context.Context, id commonids.UserAssignedIdentityId, input IdentityUpdate) (result UserAssignedIdentitiesUpdateOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUserAssignedIdentitiesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUserAssignedIdentitiesUpdate prepares the UserAssignedIdentitiesUpdate request.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesUpdate(ctx context.Context, id commonids.UserAssignedIdentityId, input IdentityUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUserAssignedIdentitiesUpdate handles the response to the UserAssignedIdentitiesUpdate request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForUserAssignedIdentitiesUpdate(resp *http.Response) (result UserAssignedIdentitiesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
