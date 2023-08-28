package storageaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HierarchicalNamespaceMigrationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type HierarchicalNamespaceMigrationOperationOptions struct {
	RequestType *string
}

func DefaultHierarchicalNamespaceMigrationOperationOptions() HierarchicalNamespaceMigrationOperationOptions {
	return HierarchicalNamespaceMigrationOperationOptions{}
}

func (o HierarchicalNamespaceMigrationOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o HierarchicalNamespaceMigrationOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.RequestType != nil {
		out["requestType"] = *o.RequestType
	}

	return out
}

// HierarchicalNamespaceMigration ...
func (c StorageAccountsClient) HierarchicalNamespaceMigration(ctx context.Context, id commonids.StorageAccountId, options HierarchicalNamespaceMigrationOperationOptions) (result HierarchicalNamespaceMigrationOperationResponse, err error) {
	req, err := c.preparerForHierarchicalNamespaceMigration(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "HierarchicalNamespaceMigration", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForHierarchicalNamespaceMigration(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "HierarchicalNamespaceMigration", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// HierarchicalNamespaceMigrationThenPoll performs HierarchicalNamespaceMigration then polls until it's completed
func (c StorageAccountsClient) HierarchicalNamespaceMigrationThenPoll(ctx context.Context, id commonids.StorageAccountId, options HierarchicalNamespaceMigrationOperationOptions) error {
	result, err := c.HierarchicalNamespaceMigration(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing HierarchicalNamespaceMigration: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after HierarchicalNamespaceMigration: %+v", err)
	}

	return nil
}

// preparerForHierarchicalNamespaceMigration prepares the HierarchicalNamespaceMigration request.
func (c StorageAccountsClient) preparerForHierarchicalNamespaceMigration(ctx context.Context, id commonids.StorageAccountId, options HierarchicalNamespaceMigrationOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/hnsonmigration", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForHierarchicalNamespaceMigration sends the HierarchicalNamespaceMigration request. The method will close the
// http.Response Body if it receives an error.
func (c StorageAccountsClient) senderForHierarchicalNamespaceMigration(ctx context.Context, req *http.Request) (future HierarchicalNamespaceMigrationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
