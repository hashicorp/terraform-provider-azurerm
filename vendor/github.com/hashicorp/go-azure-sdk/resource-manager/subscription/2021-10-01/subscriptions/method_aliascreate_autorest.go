package subscriptions

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

type AliasCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AliasCreate ...
func (c SubscriptionsClient) AliasCreate(ctx context.Context, id AliasId, input PutAliasRequest) (result AliasCreateOperationResponse, err error) {
	req, err := c.preparerForAliasCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAliasCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AliasCreateThenPoll performs AliasCreate then polls until it's completed
func (c SubscriptionsClient) AliasCreateThenPoll(ctx context.Context, id AliasId, input PutAliasRequest) error {
	result, err := c.AliasCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AliasCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AliasCreate: %+v", err)
	}

	return nil
}

// preparerForAliasCreate prepares the AliasCreate request.
func (c SubscriptionsClient) preparerForAliasCreate(ctx context.Context, id AliasId, input PutAliasRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAliasCreate sends the AliasCreate request. The method will close the
// http.Response Body if it receives an error.
func (c SubscriptionsClient) senderForAliasCreate(ctx context.Context, req *http.Request) (future AliasCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
