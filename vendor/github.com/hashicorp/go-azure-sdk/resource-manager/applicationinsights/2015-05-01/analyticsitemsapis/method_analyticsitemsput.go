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

type AnalyticsItemsPutOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApplicationInsightsComponentAnalyticsItem
}

type AnalyticsItemsPutOperationOptions struct {
	OverrideItem *bool
}

func DefaultAnalyticsItemsPutOperationOptions() AnalyticsItemsPutOperationOptions {
	return AnalyticsItemsPutOperationOptions{}
}

func (o AnalyticsItemsPutOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyticsItemsPutOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyticsItemsPutOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.OverrideItem != nil {
		out.Append("overrideItem", fmt.Sprintf("%v", *o.OverrideItem))
	}
	return &out
}

// AnalyticsItemsPut ...
func (c AnalyticsItemsAPIsClient) AnalyticsItemsPut(ctx context.Context, id ProviderComponentId, input ApplicationInsightsComponentAnalyticsItem, options AnalyticsItemsPutOperationOptions) (result AnalyticsItemsPutOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/item", id.ID()),
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

	var model ApplicationInsightsComponentAnalyticsItem
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
