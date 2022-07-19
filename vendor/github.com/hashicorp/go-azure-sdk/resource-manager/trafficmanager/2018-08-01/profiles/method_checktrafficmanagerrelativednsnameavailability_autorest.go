package profiles

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckTrafficManagerRelativeDnsNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *TrafficManagerNameAvailability
}

// CheckTrafficManagerRelativeDnsNameAvailability ...
func (c ProfilesClient) CheckTrafficManagerRelativeDnsNameAvailability(ctx context.Context, input CheckTrafficManagerRelativeDnsNameAvailabilityParameters) (result CheckTrafficManagerRelativeDnsNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckTrafficManagerRelativeDnsNameAvailability(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "CheckTrafficManagerRelativeDnsNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "CheckTrafficManagerRelativeDnsNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckTrafficManagerRelativeDnsNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "CheckTrafficManagerRelativeDnsNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckTrafficManagerRelativeDnsNameAvailability prepares the CheckTrafficManagerRelativeDnsNameAvailability request.
func (c ProfilesClient) preparerForCheckTrafficManagerRelativeDnsNameAvailability(ctx context.Context, input CheckTrafficManagerRelativeDnsNameAvailabilityParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Network/checkTrafficManagerNameAvailability"),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckTrafficManagerRelativeDnsNameAvailability handles the response to the CheckTrafficManagerRelativeDnsNameAvailability request. The method always
// closes the http.Response Body.
func (c ProfilesClient) responderForCheckTrafficManagerRelativeDnsNameAvailability(resp *http.Response) (result CheckTrafficManagerRelativeDnsNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
