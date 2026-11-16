package rules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySubscriptionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Rule
}

type ListBySubscriptionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Rule
}

type ListBySubscriptionsOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListBySubscriptionsOperationOptions() ListBySubscriptionsOperationOptions {
	return ListBySubscriptionsOperationOptions{}
}

func (o ListBySubscriptionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListBySubscriptionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListBySubscriptionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListBySubscriptionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySubscriptionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBySubscriptions ...
func (c RulesClient) ListBySubscriptions(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions) (result ListBySubscriptionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListBySubscriptionsCustomPager{},
		Path:          fmt.Sprintf("%s/rules", id.ID()),
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
		Values *[]Rule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySubscriptionsComplete retrieves all the results into a single object
func (c RulesClient) ListBySubscriptionsComplete(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions) (ListBySubscriptionsCompleteResult, error) {
	return c.ListBySubscriptionsCompleteMatchingPredicate(ctx, id, options, RuleOperationPredicate{})
}

// ListBySubscriptionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RulesClient) ListBySubscriptionsCompleteMatchingPredicate(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions, predicate RuleOperationPredicate) (result ListBySubscriptionsCompleteResult, err error) {
	items := make([]Rule, 0)

	resp, err := c.ListBySubscriptions(ctx, id, options)
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

	result = ListBySubscriptionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
