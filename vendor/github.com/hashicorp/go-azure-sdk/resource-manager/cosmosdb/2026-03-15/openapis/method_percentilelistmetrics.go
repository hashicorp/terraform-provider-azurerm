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

type PercentileListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PercentileMetric
}

type PercentileListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PercentileMetric
}

type PercentileListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPercentileListMetricsOperationOptions() PercentileListMetricsOperationOptions {
	return PercentileListMetricsOperationOptions{}
}

func (o PercentileListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PercentileListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PercentileListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type PercentileListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PercentileListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PercentileListMetrics ...
func (c OpenapisClient) PercentileListMetrics(ctx context.Context, id DatabaseAccountId, options PercentileListMetricsOperationOptions) (result PercentileListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &PercentileListMetricsCustomPager{},
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

// PercentileListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) PercentileListMetricsComplete(ctx context.Context, id DatabaseAccountId, options PercentileListMetricsOperationOptions) (PercentileListMetricsCompleteResult, error) {
	return c.PercentileListMetricsCompleteMatchingPredicate(ctx, id, options, PercentileMetricOperationPredicate{})
}

// PercentileListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) PercentileListMetricsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, options PercentileListMetricsOperationOptions, predicate PercentileMetricOperationPredicate) (result PercentileListMetricsCompleteResult, err error) {
	items := make([]PercentileMetric, 0)

	resp, err := c.PercentileListMetrics(ctx, id, options)
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

	result = PercentileListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
