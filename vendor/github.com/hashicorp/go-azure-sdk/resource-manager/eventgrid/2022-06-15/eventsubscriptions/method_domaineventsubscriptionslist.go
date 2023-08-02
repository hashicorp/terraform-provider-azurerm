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

type DomainEventSubscriptionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type DomainEventSubscriptionsListCompleteResult struct {
	Items []EventSubscription
}

type DomainEventSubscriptionsListOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultDomainEventSubscriptionsListOperationOptions() DomainEventSubscriptionsListOperationOptions {
	return DomainEventSubscriptionsListOperationOptions{}
}

func (o DomainEventSubscriptionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DomainEventSubscriptionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o DomainEventSubscriptionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// DomainEventSubscriptionsList ...
func (c EventSubscriptionsClient) DomainEventSubscriptionsList(ctx context.Context, id DomainId, options DomainEventSubscriptionsListOperationOptions) (result DomainEventSubscriptionsListOperationResponse, err error) {
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

// DomainEventSubscriptionsListComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) DomainEventSubscriptionsListComplete(ctx context.Context, id DomainId, options DomainEventSubscriptionsListOperationOptions) (DomainEventSubscriptionsListCompleteResult, error) {
	return c.DomainEventSubscriptionsListCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// DomainEventSubscriptionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) DomainEventSubscriptionsListCompleteMatchingPredicate(ctx context.Context, id DomainId, options DomainEventSubscriptionsListOperationOptions, predicate EventSubscriptionOperationPredicate) (result DomainEventSubscriptionsListCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.DomainEventSubscriptionsList(ctx, id, options)
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

	result = DomainEventSubscriptionsListCompleteResult{
		Items: items,
	}
	return
}
