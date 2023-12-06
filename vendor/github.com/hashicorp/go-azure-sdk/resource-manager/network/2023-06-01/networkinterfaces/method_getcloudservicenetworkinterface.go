package networkinterfaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCloudServiceNetworkInterfaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NetworkInterface
}

type GetCloudServiceNetworkInterfaceOperationOptions struct {
	Expand *string
}

func DefaultGetCloudServiceNetworkInterfaceOperationOptions() GetCloudServiceNetworkInterfaceOperationOptions {
	return GetCloudServiceNetworkInterfaceOperationOptions{}
}

func (o GetCloudServiceNetworkInterfaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCloudServiceNetworkInterfaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetCloudServiceNetworkInterfaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// GetCloudServiceNetworkInterface ...
func (c NetworkInterfacesClient) GetCloudServiceNetworkInterface(ctx context.Context, id RoleInstanceNetworkInterfaceId, options GetCloudServiceNetworkInterfaceOperationOptions) (result GetCloudServiceNetworkInterfaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          id.ID(),
		OptionsObject: options,
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
