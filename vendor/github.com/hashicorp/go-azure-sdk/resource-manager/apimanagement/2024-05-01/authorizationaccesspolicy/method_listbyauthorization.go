package authorizationaccesspolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAuthorizationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AuthorizationAccessPolicyContract
}

type ListByAuthorizationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AuthorizationAccessPolicyContract
}

type ListByAuthorizationOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultListByAuthorizationOperationOptions() ListByAuthorizationOperationOptions {
	return ListByAuthorizationOperationOptions{}
}

func (o ListByAuthorizationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByAuthorizationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByAuthorizationOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByAuthorizationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAuthorizationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAuthorization ...
func (c AuthorizationAccessPolicyClient) ListByAuthorization(ctx context.Context, id AuthorizationId, options ListByAuthorizationOperationOptions) (result ListByAuthorizationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByAuthorizationCustomPager{},
		Path:          fmt.Sprintf("%s/accessPolicies", id.ID()),
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
		Values *[]AuthorizationAccessPolicyContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAuthorizationComplete retrieves all the results into a single object
func (c AuthorizationAccessPolicyClient) ListByAuthorizationComplete(ctx context.Context, id AuthorizationId, options ListByAuthorizationOperationOptions) (ListByAuthorizationCompleteResult, error) {
	return c.ListByAuthorizationCompleteMatchingPredicate(ctx, id, options, AuthorizationAccessPolicyContractOperationPredicate{})
}

// ListByAuthorizationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AuthorizationAccessPolicyClient) ListByAuthorizationCompleteMatchingPredicate(ctx context.Context, id AuthorizationId, options ListByAuthorizationOperationOptions, predicate AuthorizationAccessPolicyContractOperationPredicate) (result ListByAuthorizationCompleteResult, err error) {
	items := make([]AuthorizationAccessPolicyContract, 0)

	resp, err := c.ListByAuthorization(ctx, id, options)
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

	result = ListByAuthorizationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
