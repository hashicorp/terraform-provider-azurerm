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

type PercentileSourceTargetListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PercentileMetric
}

type PercentileSourceTargetListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PercentileMetric
}

type PercentileSourceTargetListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileSourceTargetListMetricsOperationOptions() PercentileSourceTargetListMetricsOperationOptions {
	return PercentileSourceTargetListMetricsOperationOptions{}
}

func (o PercentileSourceTargetListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PercentileSourceTargetListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PercentileSourceTargetListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type PercentileSourceTargetListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PercentileSourceTargetListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PercentileSourceTargetListMetrics ...
func (c OpenapisClient) PercentileSourceTargetListMetrics(ctx context.Context, id SourceRegionTargetRegionId, options PercentileSourceTargetListMetricsOperationOptions) (result PercentileSourceTargetListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &PercentileSourceTargetListMetricsCustomPager{},
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

// PercentileSourceTargetListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) PercentileSourceTargetListMetricsComplete(ctx context.Context, id SourceRegionTargetRegionId, options PercentileSourceTargetListMetricsOperationOptions) (PercentileSourceTargetListMetricsCompleteResult, error) {
	return c.PercentileSourceTargetListMetricsCompleteMatchingPredicate(ctx, id, options, PercentileMetricOperationPredicate{})
}

// PercentileSourceTargetListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) PercentileSourceTargetListMetricsCompleteMatchingPredicate(ctx context.Context, id SourceRegionTargetRegionId, options PercentileSourceTargetListMetricsOperationOptions, predicate PercentileMetricOperationPredicate) (result PercentileSourceTargetListMetricsCompleteResult, err error) {
	items := make([]PercentileMetric, 0)

	resp, err := c.PercentileSourceTargetListMetrics(ctx, id, options)
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

	result = PercentileSourceTargetListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
