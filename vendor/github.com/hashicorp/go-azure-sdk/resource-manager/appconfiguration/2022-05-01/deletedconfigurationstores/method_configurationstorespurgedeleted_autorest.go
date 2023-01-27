package deletedconfigurationstores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStoresPurgeDeletedOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConfigurationStoresPurgeDeleted ...
func (c DeletedConfigurationStoresClient) ConfigurationStoresPurgeDeleted(ctx context.Context, id DeletedConfigurationStoreId) (result ConfigurationStoresPurgeDeletedOperationResponse, err error) {
	req, err := c.preparerForConfigurationStoresPurgeDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresPurgeDeleted", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConfigurationStoresPurgeDeleted(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresPurgeDeleted", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConfigurationStoresPurgeDeletedThenPoll performs ConfigurationStoresPurgeDeleted then polls until it's completed
func (c DeletedConfigurationStoresClient) ConfigurationStoresPurgeDeletedThenPoll(ctx context.Context, id DeletedConfigurationStoreId) error {
	result, err := c.ConfigurationStoresPurgeDeleted(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ConfigurationStoresPurgeDeleted: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConfigurationStoresPurgeDeleted: %+v", err)
	}

	return nil
}

// preparerForConfigurationStoresPurgeDeleted prepares the ConfigurationStoresPurgeDeleted request.
func (c DeletedConfigurationStoresClient) preparerForConfigurationStoresPurgeDeleted(ctx context.Context, id DeletedConfigurationStoreId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/purge", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConfigurationStoresPurgeDeleted sends the ConfigurationStoresPurgeDeleted request. The method will close the
// http.Response Body if it receives an error.
func (c DeletedConfigurationStoresClient) senderForConfigurationStoresPurgeDeleted(ctx context.Context, req *http.Request) (future ConfigurationStoresPurgeDeletedOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
