package virtualmachines

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

type PowerOffOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type PowerOffOperationOptions struct {
	SkipShutdown *bool
}

func DefaultPowerOffOperationOptions() PowerOffOperationOptions {
	return PowerOffOperationOptions{}
}

func (o PowerOffOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o PowerOffOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.SkipShutdown != nil {
		out["skipShutdown"] = *o.SkipShutdown
	}

	return out
}

// PowerOff ...
func (c VirtualMachinesClient) PowerOff(ctx context.Context, id VirtualMachineId, options PowerOffOperationOptions) (result PowerOffOperationResponse, err error) {
	req, err := c.preparerForPowerOff(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "PowerOff", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPowerOff(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "PowerOff", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PowerOffThenPoll performs PowerOff then polls until it's completed
func (c VirtualMachinesClient) PowerOffThenPoll(ctx context.Context, id VirtualMachineId, options PowerOffOperationOptions) error {
	result, err := c.PowerOff(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing PowerOff: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PowerOff: %+v", err)
	}

	return nil
}

// preparerForPowerOff prepares the PowerOff request.
func (c VirtualMachinesClient) preparerForPowerOff(ctx context.Context, id VirtualMachineId, options PowerOffOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/powerOff", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForPowerOff sends the PowerOff request. The method will close the
// http.Response Body if it receives an error.
func (c VirtualMachinesClient) senderForPowerOff(ctx context.Context, req *http.Request) (future PowerOffOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
