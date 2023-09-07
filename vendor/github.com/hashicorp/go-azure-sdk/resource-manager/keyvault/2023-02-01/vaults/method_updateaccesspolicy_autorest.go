package vaults

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateAccessPolicyOperationResponse struct {
	HttpResponse *http.Response
	Model        *VaultAccessPolicyParameters
}

// UpdateAccessPolicy ...
func (c VaultsClient) UpdateAccessPolicy(ctx context.Context, id OperationKindId, input VaultAccessPolicyParameters) (result UpdateAccessPolicyOperationResponse, err error) {
	req, err := c.preparerForUpdateAccessPolicy(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "UpdateAccessPolicy", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "UpdateAccessPolicy", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdateAccessPolicy(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "UpdateAccessPolicy", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdateAccessPolicy prepares the UpdateAccessPolicy request.
func (c VaultsClient) preparerForUpdateAccessPolicy(ctx context.Context, id OperationKindId, input VaultAccessPolicyParameters) (*http.Request, error) {
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

// responderForUpdateAccessPolicy handles the response to the UpdateAccessPolicy request. The method always
// closes the http.Response Body.
func (c VaultsClient) responderForUpdateAccessPolicy(resp *http.Response) (result UpdateAccessPolicyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
