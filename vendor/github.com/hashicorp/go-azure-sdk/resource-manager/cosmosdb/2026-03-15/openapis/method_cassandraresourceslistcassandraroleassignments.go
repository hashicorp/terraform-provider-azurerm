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

type CassandraResourcesListCassandraRoleAssignmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CassandraRoleAssignmentResource
}

type CassandraResourcesListCassandraRoleAssignmentsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CassandraRoleAssignmentResource
}

type CassandraResourcesListCassandraRoleAssignmentsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CassandraResourcesListCassandraRoleAssignmentsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CassandraResourcesListCassandraRoleAssignments ...
func (c OpenapisClient) CassandraResourcesListCassandraRoleAssignments(ctx context.Context, id DatabaseAccountId) (result CassandraResourcesListCassandraRoleAssignmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CassandraResourcesListCassandraRoleAssignmentsCustomPager{},
		Path:       fmt.Sprintf("%s/cassandraRoleAssignments", id.ID()),
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
		Values *[]CassandraRoleAssignmentResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CassandraResourcesListCassandraRoleAssignmentsComplete retrieves all the results into a single object
func (c OpenapisClient) CassandraResourcesListCassandraRoleAssignmentsComplete(ctx context.Context, id DatabaseAccountId) (CassandraResourcesListCassandraRoleAssignmentsCompleteResult, error) {
	return c.CassandraResourcesListCassandraRoleAssignmentsCompleteMatchingPredicate(ctx, id, CassandraRoleAssignmentResourceOperationPredicate{})
}

// CassandraResourcesListCassandraRoleAssignmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CassandraResourcesListCassandraRoleAssignmentsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate CassandraRoleAssignmentResourceOperationPredicate) (result CassandraResourcesListCassandraRoleAssignmentsCompleteResult, err error) {
	items := make([]CassandraRoleAssignmentResource, 0)

	resp, err := c.CassandraResourcesListCassandraRoleAssignments(ctx, id)
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

	result = CassandraResourcesListCassandraRoleAssignmentsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
