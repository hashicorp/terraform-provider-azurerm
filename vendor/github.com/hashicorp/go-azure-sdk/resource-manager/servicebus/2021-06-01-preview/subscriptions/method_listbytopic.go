package subscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByTopicOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SBSubscription
}

type ListByTopicCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SBSubscription
}

type ListByTopicOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByTopicOperationOptions() ListByTopicOperationOptions {
	return ListByTopicOperationOptions{}
}

func (o ListByTopicOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByTopicOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByTopicOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByTopicCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByTopicCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByTopic ...
func (c SubscriptionsClient) ListByTopic(ctx context.Context, id TopicId, options ListByTopicOperationOptions) (result ListByTopicOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByTopicCustomPager{},
		Path:          fmt.Sprintf("%s/subscriptions", id.ID()),
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
		Values *[]SBSubscription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByTopicComplete retrieves all the results into a single object
func (c SubscriptionsClient) ListByTopicComplete(ctx context.Context, id TopicId, options ListByTopicOperationOptions) (ListByTopicCompleteResult, error) {
	return c.ListByTopicCompleteMatchingPredicate(ctx, id, options, SBSubscriptionOperationPredicate{})
}

// ListByTopicCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SubscriptionsClient) ListByTopicCompleteMatchingPredicate(ctx context.Context, id TopicId, options ListByTopicOperationOptions, predicate SBSubscriptionOperationPredicate) (result ListByTopicCompleteResult, err error) {
	items := make([]SBSubscription, 0)

	resp, err := c.ListByTopic(ctx, id, options)
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

	result = ListByTopicCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
