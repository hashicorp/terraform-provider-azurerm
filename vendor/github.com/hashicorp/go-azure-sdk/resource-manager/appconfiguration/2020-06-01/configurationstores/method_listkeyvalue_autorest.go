package configurationstores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeyValueOperationResponse struct {
	HttpResponse *http.Response
	Model        *KeyValue
}

// ListKeyValue ...
func (c ConfigurationStoresClient) ListKeyValue(ctx context.Context, id ConfigurationStoreId, input ListKeyValueParameters) (result ListKeyValueOperationResponse, err error) {
	req, err := c.preparerForListKeyValue(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeyValue", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeyValue", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListKeyValue(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeyValue", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListKeyValue prepares the ListKeyValue request.
func (c ConfigurationStoresClient) preparerForListKeyValue(ctx context.Context, id ConfigurationStoreId, input ListKeyValueParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeyValue", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListKeyValue handles the response to the ListKeyValue request. The method always
// closes the http.Response Body.
func (c ConfigurationStoresClient) responderForListKeyValue(resp *http.Response) (result ListKeyValueOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
