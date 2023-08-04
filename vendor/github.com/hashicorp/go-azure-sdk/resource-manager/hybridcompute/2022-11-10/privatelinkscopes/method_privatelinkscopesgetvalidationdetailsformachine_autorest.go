package privatelinkscopes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesGetValidationDetailsForMachineOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkScopeValidationDetails
}

// PrivateLinkScopesGetValidationDetailsForMachine ...
func (c PrivateLinkScopesClient) PrivateLinkScopesGetValidationDetailsForMachine(ctx context.Context, id MachineId) (result PrivateLinkScopesGetValidationDetailsForMachineOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesGetValidationDetailsForMachine(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetailsForMachine", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetailsForMachine", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkScopesGetValidationDetailsForMachine(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesGetValidationDetailsForMachine", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkScopesGetValidationDetailsForMachine prepares the PrivateLinkScopesGetValidationDetailsForMachine request.
func (c PrivateLinkScopesClient) preparerForPrivateLinkScopesGetValidationDetailsForMachine(ctx context.Context, id MachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkScopes/current", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateLinkScopesGetValidationDetailsForMachine handles the response to the PrivateLinkScopesGetValidationDetailsForMachine request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesClient) responderForPrivateLinkScopesGetValidationDetailsForMachine(resp *http.Response) (result PrivateLinkScopesGetValidationDetailsForMachineOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
