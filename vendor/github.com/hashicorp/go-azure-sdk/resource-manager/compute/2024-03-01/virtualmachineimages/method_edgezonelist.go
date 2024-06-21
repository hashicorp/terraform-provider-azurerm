package virtualmachineimages

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeZoneListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualMachineImageResource
}

type EdgeZoneListOperationOptions struct {
	Expand  *string
	Orderby *string
	Top     *int64
}

func DefaultEdgeZoneListOperationOptions() EdgeZoneListOperationOptions {
	return EdgeZoneListOperationOptions{}
}

func (o EdgeZoneListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o EdgeZoneListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o EdgeZoneListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// EdgeZoneList ...
func (c VirtualMachineImagesClient) EdgeZoneList(ctx context.Context, id OfferSkuId, options EdgeZoneListOperationOptions) (result EdgeZoneListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/versions", id.ID()),
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

	var model []VirtualMachineImageResource
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
