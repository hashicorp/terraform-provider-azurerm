package cognitiveservicesaccounts

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

type CheckDomainAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *DomainAvailability
}

// CheckDomainAvailability ...
func (c CognitiveServicesAccountsClient) CheckDomainAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckDomainAvailabilityParameter) (result CheckDomainAvailabilityOperationResponse, err error) {
	req, err := c.preparerForCheckDomainAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckDomainAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckDomainAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckDomainAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "CheckDomainAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckDomainAvailability prepares the CheckDomainAvailability request.
func (c CognitiveServicesAccountsClient) preparerForCheckDomainAvailability(ctx context.Context, id commonids.SubscriptionId, input CheckDomainAvailabilityParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/checkDomainAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckDomainAvailability handles the response to the CheckDomainAvailability request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForCheckDomainAvailability(resp *http.Response) (result CheckDomainAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
