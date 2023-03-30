package redis

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

type LinkedServerCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkedServerCreate ...
func (c RedisClient) LinkedServerCreate(ctx context.Context, id LinkedServerId, input RedisLinkedServerCreateParameters) (result LinkedServerCreateOperationResponse, err error) {
	req, err := c.preparerForLinkedServerCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkedServerCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkedServerCreateThenPoll performs LinkedServerCreate then polls until it's completed
func (c RedisClient) LinkedServerCreateThenPoll(ctx context.Context, id LinkedServerId, input RedisLinkedServerCreateParameters) error {
	result, err := c.LinkedServerCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing LinkedServerCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkedServerCreate: %+v", err)
	}

	return nil
}

// preparerForLinkedServerCreate prepares the LinkedServerCreate request.
func (c RedisClient) preparerForLinkedServerCreate(ctx context.Context, id LinkedServerId, input RedisLinkedServerCreateParameters) (*http.Request, error) {
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

// senderForLinkedServerCreate sends the LinkedServerCreate request. The method will close the
// http.Response Body if it receives an error.
func (c RedisClient) senderForLinkedServerCreate(ctx context.Context, req *http.Request) (future LinkedServerCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
