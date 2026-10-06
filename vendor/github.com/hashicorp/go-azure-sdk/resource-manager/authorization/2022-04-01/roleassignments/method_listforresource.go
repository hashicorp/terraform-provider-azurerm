package roleassignments

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

type ListForResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RoleAssignment
}

type ListForResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RoleAssignment
}

type ListForResourceOperationOptions struct {
	Filter   *string
	TenantId *string
}

func DefaultListForResourceOperationOptions() ListForResourceOperationOptions {
	return ListForResourceOperationOptions{}
}

func (o ListForResourceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListForResourceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListForResourceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.TenantId != nil {
		out.Append("tenantId", fmt.Sprintf("%v", *o.TenantId))
	}
	return &out
}

type ListForResourceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListForResourceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListForResource ...
func (c RoleAssignmentsClient) ListForResource(ctx context.Context, id commonids.ScopeId, options ListForResourceOperationOptions) (result ListForResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListForResourceCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments", id.ID()),
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
		Values *[]RoleAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForResourceComplete retrieves all the results into a single object
func (c RoleAssignmentsClient) ListForResourceComplete(ctx context.Context, id commonids.ScopeId, options ListForResourceOperationOptions) (ListForResourceCompleteResult, error) {
	return c.ListForResourceCompleteMatchingPredicate(ctx, id, options, RoleAssignmentOperationPredicate{})
}

// ListForResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RoleAssignmentsClient) ListForResourceCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListForResourceOperationOptions, predicate RoleAssignmentOperationPredicate) (result ListForResourceCompleteResult, err error) {
	items := make([]RoleAssignment, 0)

	resp, err := c.ListForResource(ctx, id, options)
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

	result = ListForResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
