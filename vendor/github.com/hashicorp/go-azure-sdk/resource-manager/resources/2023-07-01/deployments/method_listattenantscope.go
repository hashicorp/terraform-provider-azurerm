package deployments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAtTenantScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentExtended
}

type ListAtTenantScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentExtended
}

type ListAtTenantScopeOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListAtTenantScopeOperationOptions() ListAtTenantScopeOperationOptions {
	return ListAtTenantScopeOperationOptions{}
}

func (o ListAtTenantScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAtTenantScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAtTenantScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListAtTenantScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAtTenantScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAtTenantScope ...
func (c DeploymentsClient) ListAtTenantScope(ctx context.Context, options ListAtTenantScopeOperationOptions) (result ListAtTenantScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAtTenantScopeCustomPager{},
		Path:          "/providers/Microsoft.Resources/deployments",
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

// ListAtTenantScopeComplete retrieves all the results into a single object
func (c DeploymentsClient) ListAtTenantScopeComplete(ctx context.Context, options ListAtTenantScopeOperationOptions) (ListAtTenantScopeCompleteResult, error) {
	return c.ListAtTenantScopeCompleteMatchingPredicate(ctx, options, DeploymentExtendedOperationPredicate{})
}

// ListAtTenantScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentsClient) ListAtTenantScopeCompleteMatchingPredicate(ctx context.Context, options ListAtTenantScopeOperationOptions, predicate DeploymentExtendedOperationPredicate) (result ListAtTenantScopeCompleteResult, err error) {
	items := make([]DeploymentExtended, 0)

	resp, err := c.ListAtTenantScope(ctx, options)
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

	result = ListAtTenantScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
