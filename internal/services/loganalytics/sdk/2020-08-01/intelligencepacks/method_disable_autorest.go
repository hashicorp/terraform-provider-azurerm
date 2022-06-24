package intelligencepacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisableOperationResponse struct {
	HttpResponse *http.Response
}

// Disable ...
func (c IntelligencePacksClient) Disable(ctx context.Context, id IntelligencePackId) (result DisableOperationResponse, err error) {
	req, err := c.preparerForDisable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "intelligencepacks.IntelligencePacksClient", "Disable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "intelligencepacks.IntelligencePacksClient", "Disable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDisable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "intelligencepacks.IntelligencePacksClient", "Disable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDisable prepares the Disable request.
func (c IntelligencePacksClient) preparerForDisable(ctx context.Context, id IntelligencePackId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDisable handles the response to the Disable request. The method always
// closes the http.Response Body.
func (c IntelligencePacksClient) responderForDisable(resp *http.Response) (result DisableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
