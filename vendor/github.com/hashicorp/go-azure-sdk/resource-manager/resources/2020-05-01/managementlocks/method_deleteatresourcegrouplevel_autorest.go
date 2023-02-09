package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteAtResourceGroupLevelOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteAtResourceGroupLevel ...
func (c ManagementLocksClient) DeleteAtResourceGroupLevel(ctx context.Context, id ProviderLockId) (result DeleteAtResourceGroupLevelOperationResponse, err error) {
	req, err := c.preparerForDeleteAtResourceGroupLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceGroupLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceGroupLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteAtResourceGroupLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceGroupLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteAtResourceGroupLevel prepares the DeleteAtResourceGroupLevel request.
func (c ManagementLocksClient) preparerForDeleteAtResourceGroupLevel(ctx context.Context, id ProviderLockId) (*http.Request, error) {
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

// responderForDeleteAtResourceGroupLevel handles the response to the DeleteAtResourceGroupLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForDeleteAtResourceGroupLevel(resp *http.Response) (result DeleteAtResourceGroupLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
