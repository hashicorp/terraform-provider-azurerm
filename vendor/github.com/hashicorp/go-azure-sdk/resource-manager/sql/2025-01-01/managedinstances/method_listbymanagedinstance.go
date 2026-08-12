package managedinstances

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByManagedInstanceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TopQueries
}

type ListByManagedInstanceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TopQueries
}

type ListByManagedInstanceOperationOptions struct {
	AggregationFunction *AggregationFunctionType
	Databases           *string
	EndTime             *string
	Interval            *QueryTimeGrainType
	NumberOfQueries     *int64
	ObservationMetric   *MetricType
	StartTime           *string
}

func DefaultListByManagedInstanceOperationOptions() ListByManagedInstanceOperationOptions {
	return ListByManagedInstanceOperationOptions{}
}

func (o ListByManagedInstanceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByManagedInstanceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByManagedInstanceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AggregationFunction != nil {
		out.Append("aggregationFunction", fmt.Sprintf("%v", *o.AggregationFunction))
	}
	if o.Databases != nil {
		out.Append("databases", fmt.Sprintf("%v", *o.Databases))
	}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.Interval != nil {
		out.Append("interval", fmt.Sprintf("%v", *o.Interval))
	}
	if o.NumberOfQueries != nil {
		out.Append("numberOfQueries", fmt.Sprintf("%v", *o.NumberOfQueries))
	}
	if o.ObservationMetric != nil {
		out.Append("observationMetric", fmt.Sprintf("%v", *o.ObservationMetric))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

type ListByManagedInstanceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByManagedInstanceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByManagedInstance ...
func (c ManagedInstancesClient) ListByManagedInstance(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByManagedInstanceOperationOptions) (result ListByManagedInstanceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByManagedInstanceCustomPager{},
		Path:          fmt.Sprintf("%s/topqueries", id.ID()),
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
		Values *[]TopQueries `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByManagedInstanceComplete retrieves all the results into a single object
func (c ManagedInstancesClient) ListByManagedInstanceComplete(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByManagedInstanceOperationOptions) (ListByManagedInstanceCompleteResult, error) {
	return c.ListByManagedInstanceCompleteMatchingPredicate(ctx, id, options, TopQueriesOperationPredicate{})
}

// ListByManagedInstanceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedInstancesClient) ListByManagedInstanceCompleteMatchingPredicate(ctx context.Context, id commonids.SqlManagedInstanceId, options ListByManagedInstanceOperationOptions, predicate TopQueriesOperationPredicate) (result ListByManagedInstanceCompleteResult, err error) {
	items := make([]TopQueries, 0)

	resp, err := c.ListByManagedInstance(ctx, id, options)
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

	result = ListByManagedInstanceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
