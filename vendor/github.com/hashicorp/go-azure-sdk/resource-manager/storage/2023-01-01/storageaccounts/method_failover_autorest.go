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

type FailoverOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type FailoverOperationOptions struct {
	FailoverType *FailoverType
}

func DefaultFailoverOperationOptions() FailoverOperationOptions {
	return FailoverOperationOptions{}
}

func (o FailoverOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o FailoverOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.FailoverType != nil {
		out["failoverType"] = *o.FailoverType
	}

	return out
}

// Failover ...
func (c StorageAccountsClient) Failover(ctx context.Context, id commonids.StorageAccountId, options FailoverOperationOptions) (result FailoverOperationResponse, err error) {
	req, err := c.preparerForFailover(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "Failover", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFailover(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "Failover", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FailoverThenPoll performs Failover then polls until it's completed
func (c StorageAccountsClient) FailoverThenPoll(ctx context.Context, id commonids.StorageAccountId, options FailoverOperationOptions) error {
	result, err := c.Failover(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Failover: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Failover: %+v", err)
	}

	return nil
}

// preparerForFailover prepares the Failover request.
func (c StorageAccountsClient) preparerForFailover(ctx context.Context, id commonids.StorageAccountId, options FailoverOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/failover", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFailover sends the Failover request. The method will close the
// http.Response Body if it receives an error.
func (c StorageAccountsClient) senderForFailover(ctx context.Context, req *http.Request) (future FailoverOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
