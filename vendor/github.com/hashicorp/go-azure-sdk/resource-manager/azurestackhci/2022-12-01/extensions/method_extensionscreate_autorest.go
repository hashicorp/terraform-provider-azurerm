package extensions

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

type ExtensionsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExtensionsCreate ...
func (c ExtensionsClient) ExtensionsCreate(ctx context.Context, id ExtensionId, input Extension) (result ExtensionsCreateOperationResponse, err error) {
	req, err := c.preparerForExtensionsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtensionsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtensionsCreateThenPoll performs ExtensionsCreate then polls until it's completed
func (c ExtensionsClient) ExtensionsCreateThenPoll(ctx context.Context, id ExtensionId, input Extension) error {
	result, err := c.ExtensionsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtensionsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtensionsCreate: %+v", err)
	}

	return nil
}

// preparerForExtensionsCreate prepares the ExtensionsCreate request.
func (c ExtensionsClient) preparerForExtensionsCreate(ctx context.Context, id ExtensionId, input Extension) (*http.Request, error) {
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

// senderForExtensionsCreate sends the ExtensionsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c ExtensionsClient) senderForExtensionsCreate(ctx context.Context, req *http.Request) (future ExtensionsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
