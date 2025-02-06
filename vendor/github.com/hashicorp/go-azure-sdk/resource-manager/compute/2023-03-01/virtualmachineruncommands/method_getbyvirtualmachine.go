package virtualmachineruncommands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetByVirtualMachineOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *VirtualMachineRunCommand
}

type GetByVirtualMachineOperationOptions struct {
	Expand *string
}

func DefaultGetByVirtualMachineOperationOptions() GetByVirtualMachineOperationOptions {
	return GetByVirtualMachineOperationOptions{}
}

func (o GetByVirtualMachineOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetByVirtualMachineOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetByVirtualMachineOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// GetByVirtualMachine ...
func (c VirtualMachineRunCommandsClient) GetByVirtualMachine(ctx context.Context, id VirtualMachineRunCommandId, options GetByVirtualMachineOperationOptions) (result GetByVirtualMachineOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          id.ID(),
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

	var model VirtualMachineRunCommand
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
