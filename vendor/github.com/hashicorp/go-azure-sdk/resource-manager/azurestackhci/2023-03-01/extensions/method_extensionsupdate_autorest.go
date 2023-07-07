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

type ExtensionsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExtensionsUpdate ...
func (c ExtensionsClient) ExtensionsUpdate(ctx context.Context, id ExtensionId, input Extension) (result ExtensionsUpdateOperationResponse, err error) {
	req, err := c.preparerForExtensionsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtensionsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtensionsUpdateThenPoll performs ExtensionsUpdate then polls until it's completed
func (c ExtensionsClient) ExtensionsUpdateThenPoll(ctx context.Context, id ExtensionId, input Extension) error {
	result, err := c.ExtensionsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtensionsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtensionsUpdate: %+v", err)
	}

	return nil
}

// preparerForExtensionsUpdate prepares the ExtensionsUpdate request.
func (c ExtensionsClient) preparerForExtensionsUpdate(ctx context.Context, id ExtensionId, input Extension) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExtensionsUpdate sends the ExtensionsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ExtensionsClient) senderForExtensionsUpdate(ctx context.Context, req *http.Request) (future ExtensionsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
