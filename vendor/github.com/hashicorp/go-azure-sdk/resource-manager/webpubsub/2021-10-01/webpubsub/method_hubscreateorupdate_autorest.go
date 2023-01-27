package webpubsub

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

type HubsCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// HubsCreateOrUpdate ...
func (c WebPubSubClient) HubsCreateOrUpdate(ctx context.Context, id HubId, input WebPubSubHub) (result HubsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForHubsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForHubsCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// HubsCreateOrUpdateThenPoll performs HubsCreateOrUpdate then polls until it's completed
func (c WebPubSubClient) HubsCreateOrUpdateThenPoll(ctx context.Context, id HubId, input WebPubSubHub) error {
	result, err := c.HubsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing HubsCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after HubsCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForHubsCreateOrUpdate prepares the HubsCreateOrUpdate request.
func (c WebPubSubClient) preparerForHubsCreateOrUpdate(ctx context.Context, id HubId, input WebPubSubHub) (*http.Request, error) {
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

// senderForHubsCreateOrUpdate sends the HubsCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c WebPubSubClient) senderForHubsCreateOrUpdate(ctx context.Context, req *http.Request) (future HubsCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
