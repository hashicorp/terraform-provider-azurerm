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

type DigitalTwinsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DigitalTwinsUpdate ...
func (c DigitalTwinsInstanceClient) DigitalTwinsUpdate(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsPatchDescription) (result DigitalTwinsUpdateOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDigitalTwinsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DigitalTwinsUpdateThenPoll performs DigitalTwinsUpdate then polls until it's completed
func (c DigitalTwinsInstanceClient) DigitalTwinsUpdateThenPoll(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsPatchDescription) error {
	result, err := c.DigitalTwinsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DigitalTwinsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DigitalTwinsUpdate: %+v", err)
	}

	return nil
}

// preparerForDigitalTwinsUpdate prepares the DigitalTwinsUpdate request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsUpdate(ctx context.Context, id DigitalTwinsInstanceId, input DigitalTwinsPatchDescription) (*http.Request, error) {
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

// senderForDigitalTwinsUpdate sends the DigitalTwinsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c DigitalTwinsInstanceClient) senderForDigitalTwinsUpdate(ctx context.Context, req *http.Request) (future DigitalTwinsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
