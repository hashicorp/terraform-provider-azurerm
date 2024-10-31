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

type SystemTopicEventSubscriptionsListBySystemTopicOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type SystemTopicEventSubscriptionsListBySystemTopicCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EventSubscription
}

type SystemTopicEventSubscriptionsListBySystemTopicOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultSystemTopicEventSubscriptionsListBySystemTopicOperationOptions() SystemTopicEventSubscriptionsListBySystemTopicOperationOptions {
	return SystemTopicEventSubscriptionsListBySystemTopicOperationOptions{}
}

func (o SystemTopicEventSubscriptionsListBySystemTopicOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o SystemTopicEventSubscriptionsListBySystemTopicOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SystemTopicEventSubscriptionsListBySystemTopicOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type SystemTopicEventSubscriptionsListBySystemTopicCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SystemTopicEventSubscriptionsListBySystemTopicCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SystemTopicEventSubscriptionsListBySystemTopic ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopic(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOperationOptions) (result SystemTopicEventSubscriptionsListBySystemTopicOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &SystemTopicEventSubscriptionsListBySystemTopicCustomPager{},
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

// SystemTopicEventSubscriptionsListBySystemTopicComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopicComplete(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOperationOptions) (SystemTopicEventSubscriptionsListBySystemTopicCompleteResult, error) {
	return c.SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOperationOptions, predicate EventSubscriptionOperationPredicate) (result SystemTopicEventSubscriptionsListBySystemTopicCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.SystemTopicEventSubscriptionsListBySystemTopic(ctx, id, options)
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

	result = SystemTopicEventSubscriptionsListBySystemTopicCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
