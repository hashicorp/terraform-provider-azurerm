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

type DatabaseAccountRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Metric
}

type DatabaseAccountRegionListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Metric
}

type DatabaseAccountRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultDatabaseAccountRegionListMetricsOperationOptions() DatabaseAccountRegionListMetricsOperationOptions {
	return DatabaseAccountRegionListMetricsOperationOptions{}
}

func (o DatabaseAccountRegionListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DatabaseAccountRegionListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DatabaseAccountRegionListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type DatabaseAccountRegionListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseAccountRegionListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseAccountRegionListMetrics ...
func (c OpenapisClient) DatabaseAccountRegionListMetrics(ctx context.Context, id RegionId, options DatabaseAccountRegionListMetricsOperationOptions) (result DatabaseAccountRegionListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DatabaseAccountRegionListMetricsCustomPager{},
		Path:          fmt.Sprintf("%s/metrics", id.ID()),
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
		Values *[]Metric `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DatabaseAccountRegionListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseAccountRegionListMetricsComplete(ctx context.Context, id RegionId, options DatabaseAccountRegionListMetricsOperationOptions) (DatabaseAccountRegionListMetricsCompleteResult, error) {
	return c.DatabaseAccountRegionListMetricsCompleteMatchingPredicate(ctx, id, options, MetricOperationPredicate{})
}

// DatabaseAccountRegionListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseAccountRegionListMetricsCompleteMatchingPredicate(ctx context.Context, id RegionId, options DatabaseAccountRegionListMetricsOperationOptions, predicate MetricOperationPredicate) (result DatabaseAccountRegionListMetricsCompleteResult, err error) {
	items := make([]Metric, 0)

	resp, err := c.DatabaseAccountRegionListMetrics(ctx, id, options)
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

	result = DatabaseAccountRegionListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
