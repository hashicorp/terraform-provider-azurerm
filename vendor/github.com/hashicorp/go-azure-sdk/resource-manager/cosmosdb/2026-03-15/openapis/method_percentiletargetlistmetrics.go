package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PercentileTargetListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PercentileMetric
}

type PercentileTargetListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PercentileMetric
}

type PercentileTargetListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileTargetListMetricsOperationOptions() PercentileTargetListMetricsOperationOptions {
	return PercentileTargetListMetricsOperationOptions{}
}

func (o PercentileTargetListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PercentileTargetListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PercentileTargetListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type PercentileTargetListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PercentileTargetListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PercentileTargetListMetrics ...
func (c OpenapisClient) PercentileTargetListMetrics(ctx context.Context, id TargetRegionId, options PercentileTargetListMetricsOperationOptions) (result PercentileTargetListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &PercentileTargetListMetricsCustomPager{},
		Path:          fmt.Sprintf("%s/percentile/metrics", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]PercentileMetric `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PercentileTargetListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) PercentileTargetListMetricsComplete(ctx context.Context, id TargetRegionId, options PercentileTargetListMetricsOperationOptions) (PercentileTargetListMetricsCompleteResult, error) {
	return c.PercentileTargetListMetricsCompleteMatchingPredicate(ctx, id, options, PercentileMetricOperationPredicate{})
}

// PercentileTargetListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) PercentileTargetListMetricsCompleteMatchingPredicate(ctx context.Context, id TargetRegionId, options PercentileTargetListMetricsOperationOptions, predicate PercentileMetricOperationPredicate) (result PercentileTargetListMetricsCompleteResult, err error) {
	items := make([]PercentileMetric, 0)

	resp, err := c.PercentileTargetListMetrics(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = PercentileTargetListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
