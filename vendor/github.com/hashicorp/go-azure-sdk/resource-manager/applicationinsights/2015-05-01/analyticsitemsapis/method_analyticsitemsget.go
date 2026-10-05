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

type AnalyticsItemsGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApplicationInsightsComponentAnalyticsItem
}

type AnalyticsItemsGetOperationOptions struct {
	Id   *string
	Name *string
}

func DefaultAnalyticsItemsGetOperationOptions() AnalyticsItemsGetOperationOptions {
	return AnalyticsItemsGetOperationOptions{}
}

func (o AnalyticsItemsGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyticsItemsGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyticsItemsGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Id != nil {
		out.Append("id", fmt.Sprintf("%v", *o.Id))
	}
	if o.Name != nil {
		out.Append("name", fmt.Sprintf("%v", *o.Name))
	}
	return &out
}

// AnalyticsItemsGet ...
func (c AnalyticsItemsAPIsClient) AnalyticsItemsGet(ctx context.Context, id ProviderComponentId, options AnalyticsItemsGetOperationOptions) (result AnalyticsItemsGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
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

	var model ApplicationInsightsComponentAnalyticsItem
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
