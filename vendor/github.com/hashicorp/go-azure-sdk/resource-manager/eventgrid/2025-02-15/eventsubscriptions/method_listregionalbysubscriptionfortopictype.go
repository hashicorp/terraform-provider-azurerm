package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRegionalBySubscriptionForTopicTypeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type ListRegionalBySubscriptionForTopicTypeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EventSubscription
}

type ListRegionalBySubscriptionForTopicTypeOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListRegionalBySubscriptionForTopicTypeOperationOptions() ListRegionalBySubscriptionForTopicTypeOperationOptions {
	return ListRegionalBySubscriptionForTopicTypeOperationOptions{}
}

func (o ListRegionalBySubscriptionForTopicTypeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListRegionalBySubscriptionForTopicTypeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListRegionalBySubscriptionForTopicTypeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListRegionalBySubscriptionForTopicTypeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRegionalBySubscriptionForTopicTypeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRegionalBySubscriptionForTopicType ...
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicType(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOperationOptions) (result ListRegionalBySubscriptionForTopicTypeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListRegionalBySubscriptionForTopicTypeCustomPager{},
		Path:          fmt.Sprintf("%s/eventSubscriptions", id.ID()),
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
		Values *[]EventSubscription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRegionalBySubscriptionForTopicTypeComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicTypeComplete(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOperationOptions) (ListRegionalBySubscriptionForTopicTypeCompleteResult, error) {
	return c.ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOperationOptions, predicate EventSubscriptionOperationPredicate) (result ListRegionalBySubscriptionForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.ListRegionalBySubscriptionForTopicType(ctx, id, options)
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

	result = ListRegionalBySubscriptionForTopicTypeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
