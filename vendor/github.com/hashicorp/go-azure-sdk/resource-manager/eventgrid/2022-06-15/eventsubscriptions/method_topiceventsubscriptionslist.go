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

type TopicEventSubscriptionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type TopicEventSubscriptionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EventSubscription
}

type TopicEventSubscriptionsListOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultTopicEventSubscriptionsListOperationOptions() TopicEventSubscriptionsListOperationOptions {
	return TopicEventSubscriptionsListOperationOptions{}
}

func (o TopicEventSubscriptionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TopicEventSubscriptionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o TopicEventSubscriptionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// TopicEventSubscriptionsList ...
func (c EventSubscriptionsClient) TopicEventSubscriptionsList(ctx context.Context, id TopicId, options TopicEventSubscriptionsListOperationOptions) (result TopicEventSubscriptionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/eventSubscriptions", id.ID()),
		OptionsObject: options,
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

// TopicEventSubscriptionsListComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) TopicEventSubscriptionsListComplete(ctx context.Context, id TopicId, options TopicEventSubscriptionsListOperationOptions) (TopicEventSubscriptionsListCompleteResult, error) {
	return c.TopicEventSubscriptionsListCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// TopicEventSubscriptionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) TopicEventSubscriptionsListCompleteMatchingPredicate(ctx context.Context, id TopicId, options TopicEventSubscriptionsListOperationOptions, predicate EventSubscriptionOperationPredicate) (result TopicEventSubscriptionsListCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.TopicEventSubscriptionsList(ctx, id, options)
	if err != nil {
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

	result = TopicEventSubscriptionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
