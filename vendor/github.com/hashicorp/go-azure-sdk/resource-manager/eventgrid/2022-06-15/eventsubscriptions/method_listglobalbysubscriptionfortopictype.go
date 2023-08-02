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

type ListGlobalBySubscriptionForTopicTypeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type ListGlobalBySubscriptionForTopicTypeCompleteResult struct {
	Items []EventSubscription
}

type ListGlobalBySubscriptionForTopicTypeOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListGlobalBySubscriptionForTopicTypeOperationOptions() ListGlobalBySubscriptionForTopicTypeOperationOptions {
	return ListGlobalBySubscriptionForTopicTypeOperationOptions{}
}

func (o ListGlobalBySubscriptionForTopicTypeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListGlobalBySubscriptionForTopicTypeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListGlobalBySubscriptionForTopicTypeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListGlobalBySubscriptionForTopicType ...
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicType(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOperationOptions) (result ListGlobalBySubscriptionForTopicTypeOperationResponse, err error) {
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

// ListGlobalBySubscriptionForTopicTypeComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicTypeComplete(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOperationOptions) (ListGlobalBySubscriptionForTopicTypeCompleteResult, error) {
	return c.ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOperationOptions, predicate EventSubscriptionOperationPredicate) (result ListGlobalBySubscriptionForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.ListGlobalBySubscriptionForTopicType(ctx, id, options)
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

	result = ListGlobalBySubscriptionForTopicTypeCompleteResult{
		Items: items,
	}
	return
}
