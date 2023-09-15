package storageaccounts

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

type RegenerateKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *StorageAccountListKeysResult
}

// RegenerateKey ...
func (c StorageAccountsClient) RegenerateKey(ctx context.Context, id commonids.StorageAccountId, input StorageAccountRegenerateKeyParameters) (result RegenerateKeyOperationResponse, err error) {
	req, err := c.preparerForRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RegenerateKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegenerateKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "RegenerateKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegenerateKey prepares the RegenerateKey request.
func (c StorageAccountsClient) preparerForRegenerateKey(ctx context.Context, id commonids.StorageAccountId, input StorageAccountRegenerateKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegenerateKey handles the response to the RegenerateKey request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForRegenerateKey(resp *http.Response) (result RegenerateKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
