package protectioncontainers

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisterOperationResponse struct {
	HttpResponse *http.Response
	Model        *ProtectionContainerResource
}

// Register ...
func (c ProtectionContainersClient) Register(ctx context.Context, id ProtectionContainerId, input ProtectionContainerResource) (result RegisterOperationResponse, err error) {
	req, err := c.preparerForRegister(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Register", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Register", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegister(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Register", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegister prepares the Register request.
func (c ProtectionContainersClient) preparerForRegister(ctx context.Context, id ProtectionContainerId, input ProtectionContainerResource) (*http.Request, error) {
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

// responderForRegister handles the response to the Register request. The method always
// closes the http.Response Body.
func (c ProtectionContainersClient) responderForRegister(resp *http.Response) (result RegisterOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
