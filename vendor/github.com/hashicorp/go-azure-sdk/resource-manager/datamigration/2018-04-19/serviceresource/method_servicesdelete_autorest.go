package serviceresource

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

type ServicesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type ServicesDeleteOperationOptions struct {
	DeleteRunningTasks *bool
}

func DefaultServicesDeleteOperationOptions() ServicesDeleteOperationOptions {
	return ServicesDeleteOperationOptions{}
}

func (o ServicesDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ServicesDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.DeleteRunningTasks != nil {
		out["deleteRunningTasks"] = *o.DeleteRunningTasks
	}

	return out
}

// ServicesDelete ...
func (c ServiceResourceClient) ServicesDelete(ctx context.Context, id ServiceId, options ServicesDeleteOperationOptions) (result ServicesDeleteOperationResponse, err error) {
	req, err := c.preparerForServicesDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForServicesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ServicesDeleteThenPoll performs ServicesDelete then polls until it's completed
func (c ServiceResourceClient) ServicesDeleteThenPoll(ctx context.Context, id ServiceId, options ServicesDeleteOperationOptions) error {
	result, err := c.ServicesDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ServicesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ServicesDelete: %+v", err)
	}

	return nil
}

// preparerForServicesDelete prepares the ServicesDelete request.
func (c ServiceResourceClient) preparerForServicesDelete(ctx context.Context, id ServiceId, options ServicesDeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForServicesDelete sends the ServicesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ServiceResourceClient) senderForServicesDelete(ctx context.Context, req *http.Request) (future ServicesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
