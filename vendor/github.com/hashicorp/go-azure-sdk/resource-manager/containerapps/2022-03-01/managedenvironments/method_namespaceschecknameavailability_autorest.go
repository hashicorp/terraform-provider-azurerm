package managedenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespacesCheckNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResponse
}

// NamespacesCheckNameAvailability ...
func (c ManagedEnvironmentsClient) NamespacesCheckNameAvailability(ctx context.Context, id ManagedEnvironmentId, input CheckNameAvailabilityRequest) (result NamespacesCheckNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForNamespacesCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "NamespacesCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "NamespacesCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForNamespacesCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "NamespacesCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForNamespacesCheckNameAvailability prepares the NamespacesCheckNameAvailability request.
func (c ManagedEnvironmentsClient) preparerForNamespacesCheckNameAvailability(ctx context.Context, id ManagedEnvironmentId, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForNamespacesCheckNameAvailability handles the response to the NamespacesCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForNamespacesCheckNameAvailability(resp *http.Response) (result NamespacesCheckNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
