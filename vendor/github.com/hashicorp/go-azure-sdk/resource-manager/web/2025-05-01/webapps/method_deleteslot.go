package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteSlotOperationOptions struct {
	DeleteEmptyServerFarm *bool
	DeleteMetrics         *bool
}

func DefaultDeleteSlotOperationOptions() DeleteSlotOperationOptions {
	return DeleteSlotOperationOptions{}
}

func (o DeleteSlotOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeleteSlotOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DeleteSlotOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DeleteEmptyServerFarm != nil {
		out.Append("deleteEmptyServerFarm", fmt.Sprintf("%v", *o.DeleteEmptyServerFarm))
	}
	if o.DeleteMetrics != nil {
		out.Append("deleteMetrics", fmt.Sprintf("%v", *o.DeleteMetrics))
	}
	return &out
}

// DeleteSlot ...
func (c WebAppsClient) DeleteSlot(ctx context.Context, id SlotId, options DeleteSlotOperationOptions) (result DeleteSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
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

	return
}
