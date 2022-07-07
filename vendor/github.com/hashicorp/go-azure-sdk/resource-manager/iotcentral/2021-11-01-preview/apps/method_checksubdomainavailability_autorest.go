package apps

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

type CheckSubdomainAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *AppAvailabilityInfo
}

// CheckSubdomainAvailability ...
func (c AppsClient) CheckSubdomainAvailability(ctx context.Context, id commonids.SubscriptionId, input OperationInputs) (result CheckSubdomainAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckSubdomainAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "CheckSubdomainAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "CheckSubdomainAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckSubdomainAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "CheckSubdomainAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckSubdomainAvailability prepares the CheckSubdomainAvailability request.
func (c AppsClient) preparerForCheckSubdomainAvailability(ctx context.Context, id commonids.SubscriptionId, input OperationInputs) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.IoTCentral/checkSubdomainAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckSubdomainAvailability handles the response to the CheckSubdomainAvailability request. The method always
// closes the http.Response Body.
func (c AppsClient) responderForCheckSubdomainAvailability(resp *http.Response) (result CheckSubdomainAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
