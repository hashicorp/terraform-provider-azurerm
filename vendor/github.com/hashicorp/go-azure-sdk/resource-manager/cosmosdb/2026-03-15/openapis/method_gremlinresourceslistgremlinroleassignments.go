package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesListGremlinRoleAssignmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GremlinRoleAssignmentResource
}

type GremlinResourcesListGremlinRoleAssignmentsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GremlinRoleAssignmentResource
}

type GremlinResourcesListGremlinRoleAssignmentsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GremlinResourcesListGremlinRoleAssignmentsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GremlinResourcesListGremlinRoleAssignments ...
func (c OpenapisClient) GremlinResourcesListGremlinRoleAssignments(ctx context.Context, id DatabaseAccountId) (result GremlinResourcesListGremlinRoleAssignmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GremlinResourcesListGremlinRoleAssignmentsCustomPager{},
		Path:       fmt.Sprintf("%s/gremlinRoleAssignments", id.ID()),
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
		Values *[]GremlinRoleAssignmentResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GremlinResourcesListGremlinRoleAssignmentsComplete retrieves all the results into a single object
func (c OpenapisClient) GremlinResourcesListGremlinRoleAssignmentsComplete(ctx context.Context, id DatabaseAccountId) (GremlinResourcesListGremlinRoleAssignmentsCompleteResult, error) {
	return c.GremlinResourcesListGremlinRoleAssignmentsCompleteMatchingPredicate(ctx, id, GremlinRoleAssignmentResourceOperationPredicate{})
}

// GremlinResourcesListGremlinRoleAssignmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) GremlinResourcesListGremlinRoleAssignmentsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate GremlinRoleAssignmentResourceOperationPredicate) (result GremlinResourcesListGremlinRoleAssignmentsCompleteResult, err error) {
	items := make([]GremlinRoleAssignmentResource, 0)

	resp, err := c.GremlinResourcesListGremlinRoleAssignments(ctx, id)
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

	result = GremlinResourcesListGremlinRoleAssignmentsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
