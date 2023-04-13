package integrationaccountpartners

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListContentCallbackUrlOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowTriggerCallbackUrl
}

// ListContentCallbackUrl ...
func (c IntegrationAccountPartnersClient) ListContentCallbackUrl(ctx context.Context, id PartnerId, input GetCallbackUrlParameters) (result ListContentCallbackUrlOperationResponse, err error) {
	req, err := c.preparerForListContentCallbackUrl(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccountpartners.IntegrationAccountPartnersClient", "ListContentCallbackUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccountpartners.IntegrationAccountPartnersClient", "ListContentCallbackUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListContentCallbackUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccountpartners.IntegrationAccountPartnersClient", "ListContentCallbackUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListContentCallbackUrl prepares the ListContentCallbackUrl request.
func (c IntegrationAccountPartnersClient) preparerForListContentCallbackUrl(ctx context.Context, id PartnerId, input GetCallbackUrlParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listContentCallbackUrl", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListContentCallbackUrl handles the response to the ListContentCallbackUrl request. The method always
// closes the http.Response Body.
func (c IntegrationAccountPartnersClient) responderForListContentCallbackUrl(resp *http.Response) (result ListContentCallbackUrlOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
