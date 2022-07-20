package redisenterprise

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

type DatabasesRegenerateKeyOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DatabasesRegenerateKey ...
func (c RedisEnterpriseClient) DatabasesRegenerateKey(ctx context.Context, id DatabaseId, input RegenerateKeyParameters) (result DatabasesRegenerateKeyOperationResponse, err error) {
	req, err := c.preparerForDatabasesRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesRegenerateKey", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDatabasesRegenerateKey(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesRegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DatabasesRegenerateKeyThenPoll performs DatabasesRegenerateKey then polls until it's completed
func (c RedisEnterpriseClient) DatabasesRegenerateKeyThenPoll(ctx context.Context, id DatabaseId, input RegenerateKeyParameters) error {
	result, err := c.DatabasesRegenerateKey(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DatabasesRegenerateKey: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DatabasesRegenerateKey: %+v", err)
	}

	return nil
}

// preparerForDatabasesRegenerateKey prepares the DatabasesRegenerateKey request.
func (c RedisEnterpriseClient) preparerForDatabasesRegenerateKey(ctx context.Context, id DatabaseId, input RegenerateKeyParameters) (*http.Request, error) {
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

// senderForDatabasesRegenerateKey sends the DatabasesRegenerateKey request. The method will close the
// http.Response Body if it receives an error.
func (c RedisEnterpriseClient) senderForDatabasesRegenerateKey(ctx context.Context, req *http.Request) (future DatabasesRegenerateKeyOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
