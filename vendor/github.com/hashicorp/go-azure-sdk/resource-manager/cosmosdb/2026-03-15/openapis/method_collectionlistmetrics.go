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

type CollectionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Metric
}

type CollectionListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Metric
}

type CollectionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionListMetricsOperationOptions() CollectionListMetricsOperationOptions {
	return CollectionListMetricsOperationOptions{}
}

func (o CollectionListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type CollectionListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionListMetrics ...
func (c OpenapisClient) CollectionListMetrics(ctx context.Context, id CollectionId, options CollectionListMetricsOperationOptions) (result CollectionListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CollectionListMetricsCustomPager{},
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

// CollectionListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionListMetricsComplete(ctx context.Context, id CollectionId, options CollectionListMetricsOperationOptions) (CollectionListMetricsCompleteResult, error) {
	return c.CollectionListMetricsCompleteMatchingPredicate(ctx, id, options, MetricOperationPredicate{})
}

// CollectionListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionListMetricsCompleteMatchingPredicate(ctx context.Context, id CollectionId, options CollectionListMetricsOperationOptions, predicate MetricOperationPredicate) (result CollectionListMetricsCompleteResult, err error) {
	items := make([]Metric, 0)

	resp, err := c.CollectionListMetrics(ctx, id, options)
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

	result = CollectionListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
