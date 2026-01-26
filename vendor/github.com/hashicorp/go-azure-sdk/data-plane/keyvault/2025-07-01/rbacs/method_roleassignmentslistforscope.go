package rbacs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentsListForScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RoleAssignment
}

type RoleAssignmentsListForScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RoleAssignment
}

type RoleAssignmentsListForScopeOperationOptions struct {
	Filter *string
}

func DefaultRoleAssignmentsListForScopeOperationOptions() RoleAssignmentsListForScopeOperationOptions {
	return RoleAssignmentsListForScopeOperationOptions{}
}

func (o RoleAssignmentsListForScopeOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RoleAssignmentsListForScopeOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RoleAssignmentsListForScopeOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type RoleAssignmentsListForScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RoleAssignmentsListForScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RoleAssignmentsListForScope ...
func (c RbacsClient) RoleAssignmentsListForScope(ctx context.Context, id ScopeId, options RoleAssignmentsListForScopeOperationOptions) (result RoleAssignmentsListForScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &RoleAssignmentsListForScopeCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments", id.Path()),
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

// RoleAssignmentsListForScopeComplete retrieves all the results into a single object
func (c RbacsClient) RoleAssignmentsListForScopeComplete(ctx context.Context, id ScopeId, options RoleAssignmentsListForScopeOperationOptions) (RoleAssignmentsListForScopeCompleteResult, error) {
	return c.RoleAssignmentsListForScopeCompleteMatchingPredicate(ctx, id, options, RoleAssignmentOperationPredicate{})
}

// RoleAssignmentsListForScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RbacsClient) RoleAssignmentsListForScopeCompleteMatchingPredicate(ctx context.Context, id ScopeId, options RoleAssignmentsListForScopeOperationOptions, predicate RoleAssignmentOperationPredicate) (result RoleAssignmentsListForScopeCompleteResult, err error) {
	items := make([]RoleAssignment, 0)

	resp, err := c.RoleAssignmentsListForScope(ctx, id, options)
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

	result = RoleAssignmentsListForScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
