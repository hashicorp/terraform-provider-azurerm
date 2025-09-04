package protectioncontainers

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

type RegisterOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
	Model        *ProtectionContainerResource
}

// Register ...
func (c ProtectionContainersClient) Register(ctx context.Context, id ProtectionContainerId, input ProtectionContainerResource) (result RegisterOperationResponse, err error) {
	req, err := c.preparerForRegister(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Register", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRegister(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Register", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RegisterThenPoll performs Register then polls until it's completed
func (c ProtectionContainersClient) RegisterThenPoll(ctx context.Context, id ProtectionContainerId, input ProtectionContainerResource) error {
	result, err := c.Register(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Register: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Register: %+v", err)
	}

	return nil
}

// preparerForRegister prepares the Register request.
func (c ProtectionContainersClient) preparerForRegister(ctx context.Context, id ProtectionContainerId, input ProtectionContainerResource) (*http.Request, error) {
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

// senderForRegister sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (c ProtectionContainersClient) senderForRegister(ctx context.Context, req *http.Request) (future RegisterOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
