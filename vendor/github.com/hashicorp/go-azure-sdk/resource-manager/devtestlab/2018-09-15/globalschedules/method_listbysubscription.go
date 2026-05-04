package globalschedules

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

type ListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Schedule
}

type ListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Schedule
}

type ListBySubscriptionOperationOptions struct {
	Expand  *string
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultListBySubscriptionOperationOptions() ListBySubscriptionOperationOptions {
	return ListBySubscriptionOperationOptions{}
}

func (o ListBySubscriptionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBySubscriptionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListBySubscriptionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBySubscription ...
func (c GlobalSchedulesClient) ListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (result ListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListBySubscriptionCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.DevTestLab/schedules", id.ID()),
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
		Values *[]Schedule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySubscriptionComplete retrieves all the results into a single object
func (c GlobalSchedulesClient) ListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (ListBySubscriptionCompleteResult, error) {
	return c.ListBySubscriptionCompleteMatchingPredicate(ctx, id, options, ScheduleOperationPredicate{})
}

// ListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalSchedulesClient) ListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions, predicate ScheduleOperationPredicate) (result ListBySubscriptionCompleteResult, err error) {
	items := make([]Schedule, 0)

	resp, err := c.ListBySubscription(ctx, id, options)
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

	result = ListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
