package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAtResourceGroupLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagementLockObject
}

// GetAtResourceGroupLevel ...
func (c ManagementLocksClient) GetAtResourceGroupLevel(ctx context.Context, id ProviderLockId) (result GetAtResourceGroupLevelOperationResponse, err error) {
	req, err := c.preparerForGetAtResourceGroupLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceGroupLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceGroupLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAtResourceGroupLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceGroupLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAtResourceGroupLevel prepares the GetAtResourceGroupLevel request.
func (c ManagementLocksClient) preparerForGetAtResourceGroupLevel(ctx context.Context, id ProviderLockId) (*http.Request, error) {
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

// responderForGetAtResourceGroupLevel handles the response to the GetAtResourceGroupLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForGetAtResourceGroupLevel(resp *http.Response) (result GetAtResourceGroupLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
