package rolemanagementpolicyassignments

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

type ListForScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RoleManagementPolicyAssignment
}

type ListForScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RoleManagementPolicyAssignment
}

type ListForScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListForScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListForScope ...
func (c RoleManagementPolicyAssignmentsClient) ListForScope(ctx context.Context, id commonids.ScopeId) (result ListForScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListForScopeCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Authorization/roleManagementPolicyAssignments", id.ID()),
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
		Values *[]RoleManagementPolicyAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListForScopeComplete retrieves all the results into a single object
func (c RoleManagementPolicyAssignmentsClient) ListForScopeComplete(ctx context.Context, id commonids.ScopeId) (ListForScopeCompleteResult, error) {
	return c.ListForScopeCompleteMatchingPredicate(ctx, id, RoleManagementPolicyAssignmentOperationPredicate{})
}

// ListForScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RoleManagementPolicyAssignmentsClient) ListForScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate RoleManagementPolicyAssignmentOperationPredicate) (result ListForScopeCompleteResult, err error) {
	items := make([]RoleManagementPolicyAssignment, 0)

	resp, err := c.ListForScope(ctx, id)
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

	result = ListForScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
