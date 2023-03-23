package trustedaccess

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleBindingsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *TrustedAccessRoleBinding
}

// RoleBindingsGet ...
func (c TrustedAccessClient) RoleBindingsGet(ctx context.Context, id TrustedAccessRoleBindingId) (result RoleBindingsGetOperationResponse, err error) {
	req, err := c.preparerForRoleBindingsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRoleBindingsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRoleBindingsGet prepares the RoleBindingsGet request.
func (c TrustedAccessClient) preparerForRoleBindingsGet(ctx context.Context, id TrustedAccessRoleBindingId) (*http.Request, error) {
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

// responderForRoleBindingsGet handles the response to the RoleBindingsGet request. The method always
// closes the http.Response Body.
func (c TrustedAccessClient) responderForRoleBindingsGet(resp *http.Response) (result RoleBindingsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
