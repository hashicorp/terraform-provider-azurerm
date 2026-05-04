package virtualmachinescalesets

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

type DeleteInstancesOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteInstancesOperationOptions struct {
	ForceDeletion *bool
}

func DefaultDeleteInstancesOperationOptions() DeleteInstancesOperationOptions {
	return DeleteInstancesOperationOptions{}
}

func (o DeleteInstancesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeleteInstancesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DeleteInstancesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ForceDeletion != nil {
		out.Append("forceDeletion", fmt.Sprintf("%v", *o.ForceDeletion))
	}
	return &out
}

// DeleteInstances ...
func (c VirtualMachineScaleSetsClient) DeleteInstances(ctx context.Context, id VirtualMachineScaleSetId, input VirtualMachineScaleSetVMInstanceRequiredIDs, options DeleteInstancesOperationOptions) (result DeleteInstancesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/delete", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// DeleteInstancesThenPoll performs DeleteInstances then polls until it's completed
func (c VirtualMachineScaleSetsClient) DeleteInstancesThenPoll(ctx context.Context, id VirtualMachineScaleSetId, input VirtualMachineScaleSetVMInstanceRequiredIDs, options DeleteInstancesOperationOptions) error {
	result, err := c.DeleteInstances(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing DeleteInstances: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeleteInstances: %+v", err)
	}

	return nil
}
