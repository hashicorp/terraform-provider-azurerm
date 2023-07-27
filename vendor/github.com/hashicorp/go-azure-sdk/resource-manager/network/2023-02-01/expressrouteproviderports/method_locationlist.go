package expressrouteproviderports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ExpressRouteProviderPortListResult
}

type LocationListOperationOptions struct {
	Filter *string
}

func DefaultLocationListOperationOptions() LocationListOperationOptions {
	return LocationListOperationOptions{}
}

func (o LocationListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocationListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o LocationListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

// LocationList ...
func (c ExpressRouteProviderPortsClient) LocationList(ctx context.Context, id commonids.SubscriptionId, options LocationListOperationOptions) (result LocationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Network/expressRouteProviderPorts", id.ID()),
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
