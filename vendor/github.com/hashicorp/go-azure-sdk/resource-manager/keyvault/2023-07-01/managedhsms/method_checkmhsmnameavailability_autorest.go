package managedhsms

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

type CheckMhsmNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckMhsmNameAvailabilityResult
}

// CheckMhsmNameAvailability ...
func (c ManagedHsmsClient) CheckMhsmNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckMhsmNameAvailabilityParameters) (result CheckMhsmNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckMhsmNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedhsms.ManagedHsmsClient", "CheckMhsmNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedhsms.ManagedHsmsClient", "CheckMhsmNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckMhsmNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedhsms.ManagedHsmsClient", "CheckMhsmNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckMhsmNameAvailability prepares the CheckMhsmNameAvailability request.
func (c ManagedHsmsClient) preparerForCheckMhsmNameAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckMhsmNameAvailabilityParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.KeyVault/checkMhsmNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckMhsmNameAvailability handles the response to the CheckMhsmNameAvailability request. The method always
// closes the http.Response Body.
func (c ManagedHsmsClient) responderForCheckMhsmNameAvailability(resp *http.Response) (result CheckMhsmNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
