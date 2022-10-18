package trustedaccess

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleBindingsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// RoleBindingsDelete ...
func (c TrustedAccessClient) RoleBindingsDelete(ctx context.Context, id TrustedAccessRoleBindingId) (result RoleBindingsDeleteOperationResponse, err error) {
	req, err := c.preparerForRoleBindingsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRoleBindingsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRoleBindingsDelete prepares the RoleBindingsDelete request.
func (c TrustedAccessClient) preparerForRoleBindingsDelete(ctx context.Context, id TrustedAccessRoleBindingId) (*http.Request, error) {
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

// responderForRoleBindingsDelete handles the response to the RoleBindingsDelete request. The method always
// closes the http.Response Body.
func (c TrustedAccessClient) responderForRoleBindingsDelete(resp *http.Response) (result RoleBindingsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
