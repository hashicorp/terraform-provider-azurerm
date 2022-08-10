package managedapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedApisGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedApiDefinition
}

// ManagedApisGet ...
func (c ManagedAPIsClient) ManagedApisGet(ctx context.Context, id ManagedApiId) (result ManagedApisGetOperationResponse, err error) {
	req, err := c.preparerForManagedApisGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedApisGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedApisGet prepares the ManagedApisGet request.
func (c ManagedAPIsClient) preparerForManagedApisGet(ctx context.Context, id ManagedApiId) (*http.Request, error) {
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

// responderForManagedApisGet handles the response to the ManagedApisGet request. The method always
// closes the http.Response Body.
func (c ManagedAPIsClient) responderForManagedApisGet(resp *http.Response) (result ManagedApisGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
