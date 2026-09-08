package profiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsGetWafLogAnalyticsMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *WafMetricsResponse
}

type LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions struct {
	Actions       *[]string
	DateTimeBegin *string
	DateTimeEnd   *string
	Granularity   *WafGranularity
	GroupBy       *[]string
	Metrics       *[]string
	RuleTypes     *[]string
}

func DefaultLogAnalyticsGetWafLogAnalyticsMetricsOperationOptions() LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions {
	return LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions{}
}

func (o LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Actions != nil {
		out.Append("actions", fmt.Sprintf("%v", *o.Actions))
	}
	if o.DateTimeBegin != nil {
		out.Append("dateTimeBegin", fmt.Sprintf("%v", *o.DateTimeBegin))
	}
	if o.DateTimeEnd != nil {
		out.Append("dateTimeEnd", fmt.Sprintf("%v", *o.DateTimeEnd))
	}
	if o.Granularity != nil {
		out.Append("granularity", fmt.Sprintf("%v", *o.Granularity))
	}
	if o.GroupBy != nil {
		out.Append("groupBy", fmt.Sprintf("%v", *o.GroupBy))
	}
	if o.Metrics != nil {
		out.Append("metrics", fmt.Sprintf("%v", *o.Metrics))
	}
	if o.RuleTypes != nil {
		out.Append("ruleTypes", fmt.Sprintf("%v", *o.RuleTypes))
	}
	return &out
}

// LogAnalyticsGetWafLogAnalyticsMetrics ...
func (c ProfilesClient) LogAnalyticsGetWafLogAnalyticsMetrics(ctx context.Context, id ProfileId, options LogAnalyticsGetWafLogAnalyticsMetricsOperationOptions) (result LogAnalyticsGetWafLogAnalyticsMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/getWafLogAnalyticsMetrics", id.ID()),
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

	var model WafMetricsResponse
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
