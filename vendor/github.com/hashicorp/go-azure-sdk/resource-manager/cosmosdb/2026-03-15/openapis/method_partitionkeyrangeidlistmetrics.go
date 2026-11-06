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

type PartitionKeyRangeIdListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PartitionMetric
}

type PartitionKeyRangeIdListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PartitionMetric
}

type PartitionKeyRangeIdListMetricsOperationOptions struct {
	Filter *string
}

func DefaultPartitionKeyRangeIdListMetricsOperationOptions() PartitionKeyRangeIdListMetricsOperationOptions {
	return PartitionKeyRangeIdListMetricsOperationOptions{}
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PartitionKeyRangeIdListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type PartitionKeyRangeIdListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PartitionKeyRangeIdListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PartitionKeyRangeIdListMetrics ...
func (c OpenapisClient) PartitionKeyRangeIdListMetrics(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions) (result PartitionKeyRangeIdListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &PartitionKeyRangeIdListMetricsCustomPager{},
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
		Values *[]PartitionMetric `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PartitionKeyRangeIdListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) PartitionKeyRangeIdListMetricsComplete(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions) (PartitionKeyRangeIdListMetricsCompleteResult, error) {
	return c.PartitionKeyRangeIdListMetricsCompleteMatchingPredicate(ctx, id, options, PartitionMetricOperationPredicate{})
}

// PartitionKeyRangeIdListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) PartitionKeyRangeIdListMetricsCompleteMatchingPredicate(ctx context.Context, id PartitionKeyRangeIdId, options PartitionKeyRangeIdListMetricsOperationOptions, predicate PartitionMetricOperationPredicate) (result PartitionKeyRangeIdListMetricsCompleteResult, err error) {
	items := make([]PartitionMetric, 0)

	resp, err := c.PartitionKeyRangeIdListMetrics(ctx, id, options)
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

	result = PartitionKeyRangeIdListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
