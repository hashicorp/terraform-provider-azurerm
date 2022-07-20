package connections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfirmConsentCodeOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConfirmConsentCodeDefinition
}

// ConfirmConsentCode ...
func (c ConnectionsClient) ConfirmConsentCode(ctx context.Context, id ConnectionId, input ConfirmConsentCodeDefinition) (result ConfirmConsentCodeOperationResponse, err error) {
	req, err := c.preparerForConfirmConsentCode(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ConfirmConsentCode", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ConfirmConsentCode", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfirmConsentCode(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ConfirmConsentCode", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfirmConsentCode prepares the ConfirmConsentCode request.
func (c ConnectionsClient) preparerForConfirmConsentCode(ctx context.Context, id ConnectionId, input ConfirmConsentCodeDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/confirmConsentCode", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConfirmConsentCode handles the response to the ConfirmConsentCode request. The method always
// closes the http.Response Body.
func (c ConnectionsClient) responderForConfirmConsentCode(resp *http.Response) (result ConfirmConsentCodeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
