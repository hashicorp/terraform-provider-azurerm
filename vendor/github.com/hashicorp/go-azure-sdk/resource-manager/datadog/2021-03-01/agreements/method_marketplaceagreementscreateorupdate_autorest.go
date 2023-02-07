package agreements

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

type MarketplaceAgreementsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatadogAgreementResource
}

// MarketplaceAgreementsCreateOrUpdate ...
func (c AgreementsClient) MarketplaceAgreementsCreateOrUpdate(ctx context.Context, id commonids.SubscriptionId, input DatadogAgreementResource) (result MarketplaceAgreementsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForMarketplaceAgreementsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMarketplaceAgreementsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMarketplaceAgreementsCreateOrUpdate prepares the MarketplaceAgreementsCreateOrUpdate request.
func (c AgreementsClient) preparerForMarketplaceAgreementsCreateOrUpdate(ctx context.Context, id commonids.SubscriptionId, input DatadogAgreementResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Datadog/agreements/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMarketplaceAgreementsCreateOrUpdate handles the response to the MarketplaceAgreementsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AgreementsClient) responderForMarketplaceAgreementsCreateOrUpdate(resp *http.Response) (result MarketplaceAgreementsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
