package configurations

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

type PutOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Put ...
func (c ConfigurationsClient) Put(ctx context.Context, id ConfigurationId, input Configuration) (result PutOperationResponse, err error) {
	req, err := c.preparerForPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "Put", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPut(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "Put", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PutThenPoll performs Put then polls until it's completed
func (c ConfigurationsClient) PutThenPoll(ctx context.Context, id ConfigurationId, input Configuration) error {
	result, err := c.Put(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Put: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Put: %+v", err)
	}

	return nil
}

// preparerForPut prepares the Put request.
func (c ConfigurationsClient) preparerForPut(ctx context.Context, id ConfigurationId, input Configuration) (*http.Request, error) {
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

// senderForPut sends the Put request. The method will close the
// http.Response Body if it receives an error.
func (c ConfigurationsClient) senderForPut(ctx context.Context, req *http.Request) (future PutOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
