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

type DigitalTwinsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DigitalTwinsCreateOrUpdate ...
func (c DigitalTwinsInstanceClient) DigitalTwinsCreateOrUpdate(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsDescription) (result DigitalTwinsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDigitalTwinsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DigitalTwinsCreateOrUpdateThenPoll performs DigitalTwinsCreateOrUpdate then polls until it's completed
func (c DigitalTwinsInstanceClient) DigitalTwinsCreateOrUpdateThenPoll(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsDescription) error {
	result, err := c.DigitalTwinsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DigitalTwinsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForDigitalTwinsCreateOrUpdate prepares the DigitalTwinsCreateOrUpdate request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsCreateOrUpdate(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsDescription) (*http.Request, error) {
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

// senderForDigitalTwinsCreateOrUpdate sends the DigitalTwinsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c DigitalTwinsInstanceClient) senderForDigitalTwinsCreateOrUpdate(ctx context.Context, req *http.Request) (future DigitalTwinsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
