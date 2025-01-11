package analyticsitemsapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyticsItemsDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type AnalyticsItemsDeleteOperationOptions struct {
	Id   *string
	Name *string
}

func DefaultAnalyticsItemsDeleteOperationOptions() AnalyticsItemsDeleteOperationOptions {
	return AnalyticsItemsDeleteOperationOptions{}
}

func (o AnalyticsItemsDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyticsItemsDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyticsItemsDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Id != nil {
		out.Append("id", fmt.Sprintf("%v", *o.Id))
	}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
	}
	return &out
}

// AnalyticsItemsDelete ...
func (c AnalyticsItemsAPIsClient) AnalyticsItemsDelete(ctx context.Context, id ProviderComponentId, options AnalyticsItemsDeleteOperationOptions) (result AnalyticsItemsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/item", id.ID()),
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
