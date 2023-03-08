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

type LinkedServerDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkedServerDelete ...
func (c RedisClient) LinkedServerDelete(ctx context.Context, id LinkedServerId) (result LinkedServerDeleteOperationResponse, err error) {
	req, err := c.preparerForLinkedServerDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkedServerDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkedServerDeleteThenPoll performs LinkedServerDelete then polls until it's completed
func (c RedisClient) LinkedServerDeleteThenPoll(ctx context.Context, id LinkedServerId) error {
	result, err := c.LinkedServerDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing LinkedServerDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkedServerDelete: %+v", err)
	}

	return nil
}

// preparerForLinkedServerDelete prepares the LinkedServerDelete request.
func (c RedisClient) preparerForLinkedServerDelete(ctx context.Context, id LinkedServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLinkedServerDelete sends the LinkedServerDelete request. The method will close the
// http.Response Body if it receives an error.
func (c RedisClient) senderForLinkedServerDelete(ctx context.Context, req *http.Request) (future LinkedServerDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
