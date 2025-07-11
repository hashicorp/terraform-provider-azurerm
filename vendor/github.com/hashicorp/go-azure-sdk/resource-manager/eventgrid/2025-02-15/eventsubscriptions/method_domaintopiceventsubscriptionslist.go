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

type DomainTopicEventSubscriptionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type DomainTopicEventSubscriptionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EventSubscription
}

type DomainTopicEventSubscriptionsListOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultDomainTopicEventSubscriptionsListOperationOptions() DomainTopicEventSubscriptionsListOperationOptions {
	return DomainTopicEventSubscriptionsListOperationOptions{}
}

func (o DomainTopicEventSubscriptionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DomainTopicEventSubscriptionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DomainTopicEventSubscriptionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type DomainTopicEventSubscriptionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DomainTopicEventSubscriptionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DomainTopicEventSubscriptionsList ...
func (c EventSubscriptionsClient) DomainTopicEventSubscriptionsList(ctx context.Context, id DomainTopicId, options DomainTopicEventSubscriptionsListOperationOptions) (result DomainTopicEventSubscriptionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DomainTopicEventSubscriptionsListCustomPager{},
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

// DomainTopicEventSubscriptionsListComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) DomainTopicEventSubscriptionsListComplete(ctx context.Context, id DomainTopicId, options DomainTopicEventSubscriptionsListOperationOptions) (DomainTopicEventSubscriptionsListCompleteResult, error) {
	return c.DomainTopicEventSubscriptionsListCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// DomainTopicEventSubscriptionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) DomainTopicEventSubscriptionsListCompleteMatchingPredicate(ctx context.Context, id DomainTopicId, options DomainTopicEventSubscriptionsListOperationOptions, predicate EventSubscriptionOperationPredicate) (result DomainTopicEventSubscriptionsListCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.DomainTopicEventSubscriptionsList(ctx, id, options)
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

	result = DomainTopicEventSubscriptionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
