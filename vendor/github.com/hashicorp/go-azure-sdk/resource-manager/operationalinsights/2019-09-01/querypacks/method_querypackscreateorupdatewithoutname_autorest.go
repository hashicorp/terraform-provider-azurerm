package querypacks

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

type QueryPacksCreateOrUpdateWithoutNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPack
}

// QueryPacksCreateOrUpdateWithoutName ...
func (c QueryPacksClient) QueryPacksCreateOrUpdateWithoutName(ctx context.Context, id commonids.ResourceGroupId, input LogAnalyticsQueryPack) (result QueryPacksCreateOrUpdateWithoutNameOperationResponse, err error) {
	req, err := c.preparerForQueryPacksCreateOrUpdateWithoutName(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksCreateOrUpdateWithoutName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksCreateOrUpdateWithoutName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueryPacksCreateOrUpdateWithoutName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksCreateOrUpdateWithoutName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueryPacksCreateOrUpdateWithoutName prepares the QueryPacksCreateOrUpdateWithoutName request.
func (c QueryPacksClient) preparerForQueryPacksCreateOrUpdateWithoutName(ctx context.Context, id commonids.ResourceGroupId, input LogAnalyticsQueryPack) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.OperationalInsights/queryPacks", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForQueryPacksCreateOrUpdateWithoutName handles the response to the QueryPacksCreateOrUpdateWithoutName request. The method always
// closes the http.Response Body.
func (c QueryPacksClient) responderForQueryPacksCreateOrUpdateWithoutName(resp *http.Response) (result QueryPacksCreateOrUpdateWithoutNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
