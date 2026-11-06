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

type TableResourcesListTableRoleAssignmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TableRoleAssignmentResource
}

type TableResourcesListTableRoleAssignmentsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TableRoleAssignmentResource
}

type TableResourcesListTableRoleAssignmentsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TableResourcesListTableRoleAssignmentsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TableResourcesListTableRoleAssignments ...
func (c OpenapisClient) TableResourcesListTableRoleAssignments(ctx context.Context, id DatabaseAccountId) (result TableResourcesListTableRoleAssignmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &TableResourcesListTableRoleAssignmentsCustomPager{},
		Path:       fmt.Sprintf("%s/tableRoleAssignments", id.ID()),
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
		Values *[]TableRoleAssignmentResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TableResourcesListTableRoleAssignmentsComplete retrieves all the results into a single object
func (c OpenapisClient) TableResourcesListTableRoleAssignmentsComplete(ctx context.Context, id DatabaseAccountId) (TableResourcesListTableRoleAssignmentsCompleteResult, error) {
	return c.TableResourcesListTableRoleAssignmentsCompleteMatchingPredicate(ctx, id, TableRoleAssignmentResourceOperationPredicate{})
}

// TableResourcesListTableRoleAssignmentsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) TableResourcesListTableRoleAssignmentsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate TableRoleAssignmentResourceOperationPredicate) (result TableResourcesListTableRoleAssignmentsCompleteResult, err error) {
	items := make([]TableRoleAssignmentResource, 0)

	resp, err := c.TableResourcesListTableRoleAssignments(ctx, id)
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

	result = TableResourcesListTableRoleAssignmentsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
