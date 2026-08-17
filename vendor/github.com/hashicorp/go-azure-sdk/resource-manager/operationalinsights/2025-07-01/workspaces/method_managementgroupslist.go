package workspaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagementGroup
}

type ManagementGroupsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagementGroup
}

type ManagementGroupsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ManagementGroupsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ManagementGroupsList ...
func (c WorkspacesClient) ManagementGroupsList(ctx context.Context, id WorkspaceId) (result ManagementGroupsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ManagementGroupsListCustomPager{},
		Path:       fmt.Sprintf("%s/managementGroups", id.ID()),
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
		Values *[]ManagementGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ManagementGroupsListComplete retrieves all the results into a single object
func (c WorkspacesClient) ManagementGroupsListComplete(ctx context.Context, id WorkspaceId) (ManagementGroupsListCompleteResult, error) {
	return c.ManagementGroupsListCompleteMatchingPredicate(ctx, id, ManagementGroupOperationPredicate{})
}

// ManagementGroupsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkspacesClient) ManagementGroupsListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate ManagementGroupOperationPredicate) (result ManagementGroupsListCompleteResult, err error) {
	items := make([]ManagementGroup, 0)

	resp, err := c.ManagementGroupsList(ctx, id)
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

	result = ManagementGroupsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
