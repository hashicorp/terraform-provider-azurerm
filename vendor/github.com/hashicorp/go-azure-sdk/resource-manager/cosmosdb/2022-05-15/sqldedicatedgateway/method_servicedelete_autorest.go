package sqldedicatedgateway

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

type ServiceDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServiceDelete ...
func (c SqlDedicatedGatewayClient) ServiceDelete(ctx context.Context, id ServiceId) (result ServiceDeleteOperationResponse, err error) {
	req, err := c.preparerForServiceDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServiceDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServiceDeleteThenPoll performs ServiceDelete then polls until it's completed
func (c SqlDedicatedGatewayClient) ServiceDeleteThenPoll(ctx context.Context, id ServiceId) error {
	result, err := c.ServiceDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ServiceDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServiceDelete: %+v", err)
	}

	return nil
}

// preparerForServiceDelete prepares the ServiceDelete request.
func (c SqlDedicatedGatewayClient) preparerForServiceDelete(ctx context.Context, id ServiceId) (*http.Request, error) {
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

// senderForServiceDelete sends the ServiceDelete request. The method will close the
// http.Response Body if it receives an error.
func (c SqlDedicatedGatewayClient) senderForServiceDelete(ctx context.Context, req *http.Request) (future ServiceDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
