package virtualnetworklinks

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

type CreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type CreateOrUpdateOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOrUpdateOperationOptions() CreateOrUpdateOperationOptions {
	return CreateOrUpdateOperationOptions{}
}

func (o CreateOrUpdateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	if o.IfNoneMatch != nil {
		out["If-None-Match"] = *o.IfNoneMatch
	}

	return out
}

func (o CreateOrUpdateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// CreateOrUpdate ...
func (c VirtualNetworkLinksClient) CreateOrUpdate(ctx context.Context, id VirtualNetworkLinkId, input VirtualNetworkLink, options CreateOrUpdateOperationOptions) (result CreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworklinks.VirtualNetworkLinksClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworklinks.VirtualNetworkLinksClient", "CreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdateThenPoll performs CreateOrUpdate then polls until it's completed
func (c VirtualNetworkLinksClient) CreateOrUpdateThenPoll(ctx context.Context, id VirtualNetworkLinkId, input VirtualNetworkLink, options CreateOrUpdateOperationOptions) error {
	result, err := c.CreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForCreateOrUpdate prepares the CreateOrUpdate request.
func (c VirtualNetworkLinksClient) preparerForCreateOrUpdate(ctx context.Context, id VirtualNetworkLinkId, input VirtualNetworkLink, options CreateOrUpdateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateOrUpdate sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualNetworkLinksClient) senderForCreateOrUpdate(ctx context.Context, req *http.Request) (future CreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
