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

type DatabaseListMetricsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Metric
}

type DatabaseListMetricsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Metric
}

type DatabaseListMetricsOperationOptions struct {
	Filter *string
}

func DefaultDatabaseListMetricsOperationOptions() DatabaseListMetricsOperationOptions {
	return DatabaseListMetricsOperationOptions{}
}

func (o DatabaseListMetricsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DatabaseListMetricsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DatabaseListMetricsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type DatabaseListMetricsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseListMetricsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseListMetrics ...
func (c OpenapisClient) DatabaseListMetrics(ctx context.Context, id DatabaseId, options DatabaseListMetricsOperationOptions) (result DatabaseListMetricsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DatabaseListMetricsCustomPager{},
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

// DatabaseListMetricsComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseListMetricsComplete(ctx context.Context, id DatabaseId, options DatabaseListMetricsOperationOptions) (DatabaseListMetricsCompleteResult, error) {
	return c.DatabaseListMetricsCompleteMatchingPredicate(ctx, id, options, MetricOperationPredicate{})
}

// DatabaseListMetricsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseListMetricsCompleteMatchingPredicate(ctx context.Context, id DatabaseId, options DatabaseListMetricsOperationOptions, predicate MetricOperationPredicate) (result DatabaseListMetricsCompleteResult, err error) {
	items := make([]Metric, 0)

	resp, err := c.DatabaseListMetrics(ctx, id, options)
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

	result = DatabaseListMetricsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
