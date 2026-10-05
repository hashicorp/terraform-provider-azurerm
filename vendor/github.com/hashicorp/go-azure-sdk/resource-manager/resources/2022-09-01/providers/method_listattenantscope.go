package providers

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
	Model        *[]Provider
}

type ListAtTenantScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Provider
}

type ListAtTenantScopeOperationOptions struct {
	Expand *string
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
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
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
func (c ProvidersClient) ListAtTenantScope(ctx context.Context, options ListAtTenantScopeOperationOptions) (result ListAtTenantScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListAtTenantScopeCustomPager{},
		Path:          "/providers",
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
		Values *[]Provider `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAtTenantScopeComplete retrieves all the results into a single object
func (c ProvidersClient) ListAtTenantScopeComplete(ctx context.Context, options ListAtTenantScopeOperationOptions) (ListAtTenantScopeCompleteResult, error) {
	return c.ListAtTenantScopeCompleteMatchingPredicate(ctx, options, ProviderOperationPredicate{})
}

// ListAtTenantScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvidersClient) ListAtTenantScopeCompleteMatchingPredicate(ctx context.Context, options ListAtTenantScopeOperationOptions, predicate ProviderOperationPredicate) (result ListAtTenantScopeCompleteResult, err error) {
	items := make([]Provider, 0)

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
