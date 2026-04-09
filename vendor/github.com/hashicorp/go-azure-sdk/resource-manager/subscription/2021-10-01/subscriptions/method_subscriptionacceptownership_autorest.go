package subscriptions

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

type SubscriptionAcceptOwnershipOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SubscriptionAcceptOwnership ...
func (c SubscriptionsClient) SubscriptionAcceptOwnership(ctx context.Context, id ProviderSubscriptionId, input AcceptOwnershipRequest) (result SubscriptionAcceptOwnershipOperationResponse, err error) {
	req, err := c.preparerForSubscriptionAcceptOwnership(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionAcceptOwnership", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSubscriptionAcceptOwnership(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionAcceptOwnership", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SubscriptionAcceptOwnershipThenPoll performs SubscriptionAcceptOwnership then polls until it's completed
func (c SubscriptionsClient) SubscriptionAcceptOwnershipThenPoll(ctx context.Context, id ProviderSubscriptionId, input AcceptOwnershipRequest) error {
	result, err := c.SubscriptionAcceptOwnership(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SubscriptionAcceptOwnership: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SubscriptionAcceptOwnership: %+v", err)
	}

	return nil
}

// preparerForSubscriptionAcceptOwnership prepares the SubscriptionAcceptOwnership request.
func (c SubscriptionsClient) preparerForSubscriptionAcceptOwnership(ctx context.Context, id ProviderSubscriptionId, input AcceptOwnershipRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/acceptOwnership", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSubscriptionAcceptOwnership sends the SubscriptionAcceptOwnership request. The method will close the
// http.Response Body if it receives an error.
func (c SubscriptionsClient) senderForSubscriptionAcceptOwnership(ctx context.Context, req *http.Request) (future SubscriptionAcceptOwnershipOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
