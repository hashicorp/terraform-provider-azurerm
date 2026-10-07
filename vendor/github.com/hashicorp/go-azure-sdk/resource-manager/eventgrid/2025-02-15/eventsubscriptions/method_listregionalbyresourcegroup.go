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

type ListRegionalByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EventSubscription
}

type ListRegionalByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EventSubscription
}

type ListRegionalByResourceGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListRegionalByResourceGroupOperationOptions() ListRegionalByResourceGroupOperationOptions {
	return ListRegionalByResourceGroupOperationOptions{}
}

func (o ListRegionalByResourceGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListRegionalByResourceGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListRegionalByResourceGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListRegionalByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRegionalByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRegionalByResourceGroup ...
func (c EventSubscriptionsClient) ListRegionalByResourceGroup(ctx context.Context, id ProviderLocationId, options ListRegionalByResourceGroupOperationOptions) (result ListRegionalByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListRegionalByResourceGroupCustomPager{},
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

// ListRegionalByResourceGroupComplete retrieves all the results into a single object
func (c EventSubscriptionsClient) ListRegionalByResourceGroupComplete(ctx context.Context, id ProviderLocationId, options ListRegionalByResourceGroupOperationOptions) (ListRegionalByResourceGroupCompleteResult, error) {
	return c.ListRegionalByResourceGroupCompleteMatchingPredicate(ctx, id, options, EventSubscriptionOperationPredicate{})
}

// ListRegionalByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventSubscriptionsClient) ListRegionalByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, options ListRegionalByResourceGroupOperationOptions, predicate EventSubscriptionOperationPredicate) (result ListRegionalByResourceGroupCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	resp, err := c.ListRegionalByResourceGroup(ctx, id, options)
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

	result = ListRegionalByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
