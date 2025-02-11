package virtualmachinescalesetvms

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PowerOffOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type PowerOffOperationOptions struct {
	SkipShutdown *bool
}

func DefaultPowerOffOperationOptions() PowerOffOperationOptions {
	return PowerOffOperationOptions{}
}

func (o PowerOffOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PowerOffOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PowerOffOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SkipShutdown != nil {
		out.Append("skipShutdown", fmt.Sprintf("%v", *o.SkipShutdown))
	}
	return &out
}

// PowerOff ...
func (c VirtualMachineScaleSetVMsClient) PowerOff(ctx context.Context, id VirtualMachineScaleSetVirtualMachineId, options PowerOffOperationOptions) (result PowerOffOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/poweroff", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// PowerOffThenPoll performs PowerOff then polls until it's completed
func (c VirtualMachineScaleSetVMsClient) PowerOffThenPoll(ctx context.Context, id VirtualMachineScaleSetVirtualMachineId, options PowerOffOperationOptions) error {
	result, err := c.PowerOff(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing PowerOff: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after PowerOff: %+v", err)
	}

	return nil
}
