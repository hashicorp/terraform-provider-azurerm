package namespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckAvailabilityResult
}

// CheckAvailability ...
func (c NamespacesClient) CheckAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckAvailabilityParameters) (result CheckAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CheckAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CheckAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "CheckAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckAvailability prepares the CheckAvailability request.
func (c NamespacesClient) preparerForCheckAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckAvailabilityParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.NotificationHubs/checkNamespaceAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckAvailability handles the response to the CheckAvailability request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForCheckAvailability(resp *http.Response) (result CheckAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
