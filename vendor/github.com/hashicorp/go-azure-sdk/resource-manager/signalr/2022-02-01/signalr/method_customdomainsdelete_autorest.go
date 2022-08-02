package signalr

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

type CustomDomainsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CustomDomainsDelete ...
func (c SignalRClient) CustomDomainsDelete(ctx context.Context, id CustomDomainId) (result CustomDomainsDeleteOperationResponse, err error) {
	req, err := c.preparerForCustomDomainsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCustomDomainsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CustomDomainsDeleteThenPoll performs CustomDomainsDelete then polls until it's completed
func (c SignalRClient) CustomDomainsDeleteThenPoll(ctx context.Context, id CustomDomainId) error {
	result, err := c.CustomDomainsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CustomDomainsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CustomDomainsDelete: %+v", err)
	}

	return nil
}

// preparerForCustomDomainsDelete prepares the CustomDomainsDelete request.
func (c SignalRClient) preparerForCustomDomainsDelete(ctx context.Context, id CustomDomainId) (*http.Request, error) {
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

// senderForCustomDomainsDelete sends the CustomDomainsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForCustomDomainsDelete(ctx context.Context, req *http.Request) (future CustomDomainsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
