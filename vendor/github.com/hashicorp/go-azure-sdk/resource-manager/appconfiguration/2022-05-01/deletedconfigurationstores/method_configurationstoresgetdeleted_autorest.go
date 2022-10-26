package deletedconfigurationstores

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStoresGetDeletedOperationResponse struct {
	HttpResponse *http.Response
	Model        *DeletedConfigurationStore
}

// ConfigurationStoresGetDeleted ...
func (c DeletedConfigurationStoresClient) ConfigurationStoresGetDeleted(ctx context.Context, id DeletedConfigurationStoreId) (result ConfigurationStoresGetDeletedOperationResponse, err error) {
	req, err := c.preparerForConfigurationStoresGetDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresGetDeleted", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresGetDeleted", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfigurationStoresGetDeleted(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresGetDeleted", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfigurationStoresGetDeleted prepares the ConfigurationStoresGetDeleted request.
func (c DeletedConfigurationStoresClient) preparerForConfigurationStoresGetDeleted(ctx context.Context, id DeletedConfigurationStoreId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConfigurationStoresGetDeleted handles the response to the ConfigurationStoresGetDeleted request. The method always
// closes the http.Response Body.
func (c DeletedConfigurationStoresClient) responderForConfigurationStoresGetDeleted(resp *http.Response) (result ConfigurationStoresGetDeletedOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
