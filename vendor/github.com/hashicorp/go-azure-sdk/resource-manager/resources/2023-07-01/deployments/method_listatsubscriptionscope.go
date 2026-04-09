package deployments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAtSubscriptionScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentExtended
}

type ListAtSubscriptionScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentExtended
}

type ListAtSubscriptionScopeOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListAtSubscriptionScopeOperationOptions() ListAtSubscriptionScopeOperationOptions {
	return ListAtSubscriptionScopeOperationOptions{}
}

func (o ListAtSubscriptionScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAtSubscriptionScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAtSubscriptionScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListAtSubscriptionScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAtSubscriptionScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAtSubscriptionScope ...
func (c DeploymentsClient) ListAtSubscriptionScope(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionScopeOperationOptions) (result ListAtSubscriptionScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAtSubscriptionScopeCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Resources/deployments", id.ID()),
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
		Values *[]DeploymentExtended `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAtSubscriptionScopeComplete retrieves all the results into a single object
func (c DeploymentsClient) ListAtSubscriptionScopeComplete(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionScopeOperationOptions) (ListAtSubscriptionScopeCompleteResult, error) {
	return c.ListAtSubscriptionScopeCompleteMatchingPredicate(ctx, id, options, DeploymentExtendedOperationPredicate{})
}

// ListAtSubscriptionScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentsClient) ListAtSubscriptionScopeCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListAtSubscriptionScopeOperationOptions, predicate DeploymentExtendedOperationPredicate) (result ListAtSubscriptionScopeCompleteResult, err error) {
	items := make([]DeploymentExtended, 0)

	resp, err := c.ListAtSubscriptionScope(ctx, id, options)
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

	result = ListAtSubscriptionScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
