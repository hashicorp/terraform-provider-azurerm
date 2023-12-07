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

type ServiceCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ServiceCreate ...
func (c SqlDedicatedGatewayClient) ServiceCreate(ctx context.Context, id ServiceId, input ServiceResourceCreateUpdateParameters) (result ServiceCreateOperationResponse, err error) {
	req, err := c.preparerForServiceCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServiceCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqldedicatedgateway.SqlDedicatedGatewayClient", "ServiceCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServiceCreateThenPoll performs ServiceCreate then polls until it's completed
func (c SqlDedicatedGatewayClient) ServiceCreateThenPoll(ctx context.Context, id ServiceId, input ServiceResourceCreateUpdateParameters) error {
	result, err := c.ServiceCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ServiceCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServiceCreate: %+v", err)
	}

	return nil
}

// preparerForServiceCreate prepares the ServiceCreate request.
func (c SqlDedicatedGatewayClient) preparerForServiceCreate(ctx context.Context, id ServiceId, input ServiceResourceCreateUpdateParameters) (*http.Request, error) {
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

// senderForServiceCreate sends the ServiceCreate request. The method will close the
// http.Response Body if it receives an error.
func (c SqlDedicatedGatewayClient) senderForServiceCreate(ctx context.Context, req *http.Request) (future ServiceCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
