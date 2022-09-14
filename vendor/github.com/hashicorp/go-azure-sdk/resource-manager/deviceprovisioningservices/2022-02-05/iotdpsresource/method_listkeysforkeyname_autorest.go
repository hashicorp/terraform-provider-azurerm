package iotdpsresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeysForKeyNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedAccessSignatureAuthorizationRuleAccessRightsDescription
}

// ListKeysForKeyName ...
func (c IotDpsResourceClient) ListKeysForKeyName(ctx context.Context, id KeyId) (result ListKeysForKeyNameOperationResponse, err error) {
	req, err := c.preparerForListKeysForKeyName(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeysForKeyName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeysForKeyName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKeysForKeyName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeysForKeyName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKeysForKeyName prepares the ListKeysForKeyName request.
func (c IotDpsResourceClient) preparerForListKeysForKeyName(ctx context.Context, id KeyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listkeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKeysForKeyName handles the response to the ListKeysForKeyName request. The method always
// closes the http.Response Body.
func (c IotDpsResourceClient) responderForListKeysForKeyName(resp *http.Response) (result ListKeysForKeyNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
