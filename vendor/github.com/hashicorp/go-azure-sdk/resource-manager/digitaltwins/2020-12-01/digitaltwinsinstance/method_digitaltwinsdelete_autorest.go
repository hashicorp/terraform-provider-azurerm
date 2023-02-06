package digitaltwinsinstance

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

type DigitalTwinsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DigitalTwinsDelete ...
func (c DigitalTwinsInstanceClient) DigitalTwinsDelete(ctx context.Context, id DigitalTwinsInstanceId) (result DigitalTwinsDeleteOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDigitalTwinsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DigitalTwinsDeleteThenPoll performs DigitalTwinsDelete then polls until it's completed
func (c DigitalTwinsInstanceClient) DigitalTwinsDeleteThenPoll(ctx context.Context, id DigitalTwinsInstanceId) error {
	result, err := c.DigitalTwinsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DigitalTwinsDelete: %+v", err)
	}

	return nil
}

// preparerForDigitalTwinsDelete prepares the DigitalTwinsDelete request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsDelete(ctx context.Context, id DigitalTwinsInstanceId) (*http.Request, error) {
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

// senderForDigitalTwinsDelete sends the DigitalTwinsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c DigitalTwinsInstanceClient) senderForDigitalTwinsDelete(ctx context.Context, req *http.Request) (future DigitalTwinsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
