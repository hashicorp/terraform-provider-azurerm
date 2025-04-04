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

type NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Subscription
}

type NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Subscription
}

type NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultNamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions() NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions {
	return NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions{}
}

func (o NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type NamespaceTopicEventSubscriptionsListByNamespaceTopicCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NamespaceTopicEventSubscriptionsListByNamespaceTopicCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NamespaceTopicEventSubscriptionsListByNamespaceTopic ...
func (c EventSubscriptionsClient) NamespaceTopicEventSubscriptionsListByNamespaceTopic(ctx context.Context, id NamespaceTopicId, options NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions) (result NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &NamespaceTopicEventSubscriptionsListByNamespaceTopicCustomPager{},
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
		Values *[]Subscription `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NamespaceTopicEventSubscriptionsListByNamespaceTopicComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) NamespaceTopicEventSubscriptionsListByNamespaceTopicComplete(ctx context.Context, id NamespaceTopicId, options NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions) (NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteResult, error) {
	return c.NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteMatchingPredicate(ctx, id, options, SubscriptionOperationPredicate{})
}

// NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteMatchingPredicate(ctx context.Context, id NamespaceTopicId, options NamespaceTopicEventSubscriptionsListByNamespaceTopicOperationOptions, predicate SubscriptionOperationPredicate) (result NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteResult, err error) {
	items := make([]Subscription, 0)

	resp, err := c.NamespaceTopicEventSubscriptionsListByNamespaceTopic(ctx, id, options)
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

	result = NamespaceTopicEventSubscriptionsListByNamespaceTopicCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
