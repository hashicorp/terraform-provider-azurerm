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

type ListByDomainTopicOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type ListByDomainTopicCompleteResult struct {
	Items []EventSubscription
}

type ListByDomainTopicOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByDomainTopicOperationOptions() ListByDomainTopicOperationOptions {
	return ListByDomainTopicOperationOptions{}
}

func (o ListByDomainTopicOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByDomainTopicOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByDomainTopicOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByDomainTopic ...
func (c EventSubscriptionsClient) ListByDomainTopic(ctx context.Context, id DomainTopicId, options ListByDomainTopicOperationOptions) (result ListByDomainTopicOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.EventGrid/eventSubscriptions", id.ID()),
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

// ListByDomainTopicComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) ListByDomainTopicComplete(ctx context.Context, id DomainTopicId, options ListByDomainTopicOperationOptions) (ListByDomainTopicCompleteResult, error) {
	return c.ListByDomainTopicCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// ListByDomainTopicCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) ListByDomainTopicCompleteMatchingPredicate(ctx context.Context, id DomainTopicId, options ListByDomainTopicOperationOptions, predicate EventSubscriptionOperationPredicate) (result ListByDomainTopicCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.ListByDomainTopic(ctx, id, options)
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

	result = ListByDomainTopicCompleteResult{
		Items: items,
	}
	return
}
