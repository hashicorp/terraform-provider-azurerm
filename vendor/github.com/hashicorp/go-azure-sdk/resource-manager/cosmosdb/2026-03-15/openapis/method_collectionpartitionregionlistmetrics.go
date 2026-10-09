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

type CollectionPartitionRegionListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PartitionMetric
}

type CollectionPartitionRegionListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PartitionMetric
}

type CollectionPartitionRegionListMetricsOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionRegionListMetricsOperationOptions() CollectionPartitionRegionListMetricsOperationOptions {
	return CollectionPartitionRegionListMetricsOperationOptions{}
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionPartitionRegionListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type CollectionPartitionRegionListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionPartitionRegionListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionPartitionRegionListMetrics ...
func (c OpenapisClient) CollectionPartitionRegionListMetrics(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions) (result CollectionPartitionRegionListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CollectionPartitionRegionListMetricsCustomPager{},
		Path:          fmt.Sprintf("%s/partitions/metrics", id.ID()),
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
		Values *[]PartitionMetric `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CollectionPartitionRegionListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionPartitionRegionListMetricsComplete(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions) (CollectionPartitionRegionListMetricsCompleteResult, error) {
	return c.CollectionPartitionRegionListMetricsCompleteMatchingPredicate(ctx, id, options, PartitionMetricOperationPredicate{})
}

// CollectionPartitionRegionListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionPartitionRegionListMetricsCompleteMatchingPredicate(ctx context.Context, id DatabaseCollectionId, options CollectionPartitionRegionListMetricsOperationOptions, predicate PartitionMetricOperationPredicate) (result CollectionPartitionRegionListMetricsCompleteResult, err error) {
	items := make([]PartitionMetric, 0)

	resp, err := c.CollectionPartitionRegionListMetrics(ctx, id, options)
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

	result = CollectionPartitionRegionListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
