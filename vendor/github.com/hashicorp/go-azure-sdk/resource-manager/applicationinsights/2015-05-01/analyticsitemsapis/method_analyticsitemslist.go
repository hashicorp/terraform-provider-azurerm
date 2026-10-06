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

type AnalyticsItemsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationInsightsComponentAnalyticsItem
}

type AnalyticsItemsListOperationOptions struct {
	IncludeContent *bool
	Scope          *ItemScope
	Type           *ItemTypeParameter
}

func DefaultAnalyticsItemsListOperationOptions() AnalyticsItemsListOperationOptions {
	return AnalyticsItemsListOperationOptions{}
}

func (o AnalyticsItemsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AnalyticsItemsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AnalyticsItemsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludeContent != nil {
		out.Append("includeContent", fmt.Sprintf("%v", *o.IncludeContent))
	}
	if o.Scope != nil {
		out.Append("scope", fmt.Sprintf("%v", *o.Scope))
	}
	if o.Type != nil {
		out.Append("type", fmt.Sprintf("%v", *o.Type))
	}
	return &out
}

// AnalyticsItemsList ...
func (c AnalyticsItemsAPIsClient) AnalyticsItemsList(ctx context.Context, id ProviderComponentId, options AnalyticsItemsListOperationOptions) (result AnalyticsItemsListOperationResponse, err error) {
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

	var model []ApplicationInsightsComponentAnalyticsItem
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
