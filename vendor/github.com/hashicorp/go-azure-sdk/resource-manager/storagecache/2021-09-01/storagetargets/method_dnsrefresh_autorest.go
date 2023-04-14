package storagetargets

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

type DnsRefreshOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DnsRefresh ...
func (c StorageTargetsClient) DnsRefresh(ctx context.Context, id StorageTargetId) (result DnsRefreshOperationResponse, err error) {
	req, err := c.preparerForDnsRefresh(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "DnsRefresh", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDnsRefresh(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "DnsRefresh", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DnsRefreshThenPoll performs DnsRefresh then polls until it's completed
func (c StorageTargetsClient) DnsRefreshThenPoll(ctx context.Context, id StorageTargetId) error {
	result, err := c.DnsRefresh(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DnsRefresh: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DnsRefresh: %+v", err)
	}

	return nil
}

// preparerForDnsRefresh prepares the DnsRefresh request.
func (c StorageTargetsClient) preparerForDnsRefresh(ctx context.Context, id StorageTargetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/dnsRefresh", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDnsRefresh sends the DnsRefresh request. The method will close the
// http.Response Body if it receives an error.
func (c StorageTargetsClient) senderForDnsRefresh(ctx context.Context, req *http.Request) (future DnsRefreshOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
