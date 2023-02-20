package share

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

type ProviderShareSubscriptionsRevokeOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ProviderShareSubscriptionsRevoke ...
func (c ShareClient) ProviderShareSubscriptionsRevoke(ctx context.Context, id ProviderShareSubscriptionId) (result ProviderShareSubscriptionsRevokeOperationResponse, err error) {
	req, err := c.preparerForProviderShareSubscriptionsRevoke(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsRevoke", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForProviderShareSubscriptionsRevoke(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsRevoke", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ProviderShareSubscriptionsRevokeThenPoll performs ProviderShareSubscriptionsRevoke then polls until it's completed
func (c ShareClient) ProviderShareSubscriptionsRevokeThenPoll(ctx context.Context, id ProviderShareSubscriptionId) error {
	result, err := c.ProviderShareSubscriptionsRevoke(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ProviderShareSubscriptionsRevoke: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ProviderShareSubscriptionsRevoke: %+v", err)
	}

	return nil
}

// preparerForProviderShareSubscriptionsRevoke prepares the ProviderShareSubscriptionsRevoke request.
func (c ShareClient) preparerForProviderShareSubscriptionsRevoke(ctx context.Context, id ProviderShareSubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/revoke", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForProviderShareSubscriptionsRevoke sends the ProviderShareSubscriptionsRevoke request. The method will close the
// http.Response Body if it receives an error.
func (c ShareClient) senderForProviderShareSubscriptionsRevoke(ctx context.Context, req *http.Request) (future ProviderShareSubscriptionsRevokeOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
