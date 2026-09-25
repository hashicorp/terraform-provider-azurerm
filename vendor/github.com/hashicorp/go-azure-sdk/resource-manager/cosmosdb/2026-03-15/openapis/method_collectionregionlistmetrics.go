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

type CollectionRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Metric
}

type CollectionRegionListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Metric
}

type CollectionRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionRegionListMetricsOperationOptions() CollectionRegionListMetricsOperationOptions {
	return CollectionRegionListMetricsOperationOptions{}
}

func (o CollectionRegionListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionRegionListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionRegionListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type CollectionRegionListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionRegionListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionRegionListMetrics ...
func (c OpenapisClient) CollectionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionRegionListMetricsOperationOptions) (result CollectionRegionListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CollectionRegionListMetricsCustomPager{},
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

// CollectionRegionListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionRegionListMetricsComplete(ctx context.Context, id DatabaseCollectionId, options CollectionRegionListMetricsOperationOptions) (CollectionRegionListMetricsCompleteResult, error) {
	return c.CollectionRegionListMetricsCompleteMatchingPredicate(ctx, id, options, MetricOperationPredicate{})
}

// CollectionRegionListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionRegionListMetricsCompleteMatchingPredicate(ctx context.Context, id DatabaseCollectionId, options CollectionRegionListMetricsOperationOptions, predicate MetricOperationPredicate) (result CollectionRegionListMetricsCompleteResult, err error) {
	items := make([]Metric, 0)

	resp, err := c.CollectionRegionListMetrics(ctx, id, options)
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

	result = CollectionRegionListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
