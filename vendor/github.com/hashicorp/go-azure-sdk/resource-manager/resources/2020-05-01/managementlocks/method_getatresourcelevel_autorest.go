package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAtResourceLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagementLockObject
}

// GetAtResourceLevel ...
func (c ManagementLocksClient) GetAtResourceLevel(ctx context.Context, id ResourceLockId) (result GetAtResourceLevelOperationResponse, err error) {
	req, err := c.preparerForGetAtResourceLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAtResourceLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtResourceLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAtResourceLevel prepares the GetAtResourceLevel request.
func (c ManagementLocksClient) preparerForGetAtResourceLevel(ctx context.Context, id ResourceLockId) (*http.Request, error) {
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

// responderForGetAtResourceLevel handles the response to the GetAtResourceLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForGetAtResourceLevel(resp *http.Response) (result GetAtResourceLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
