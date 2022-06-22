package managedapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedApisListOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagedApiDefinitionCollection
}

// ManagedApisList ...
func (c ManagedAPIsClient) ManagedApisList(ctx context.Context, id LocationId) (result ManagedApisListOperationResponse, err error) {
	req, err := c.preparerForManagedApisList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForManagedApisList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedapis.ManagedAPIsClient", "ManagedApisList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForManagedApisList prepares the ManagedApisList request.
func (c ManagedAPIsClient) preparerForManagedApisList(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/managedApis", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForManagedApisList handles the response to the ManagedApisList request. The method always
// closes the http.Response Body.
func (c ManagedAPIsClient) responderForManagedApisList(resp *http.Response) (result ManagedApisListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
