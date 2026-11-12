package goldengateconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAssignedDeploymentsByParentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AssignedDeployment
}

type ListAssignedDeploymentsByParentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AssignedDeployment
}

type ListAssignedDeploymentsByParentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAssignedDeploymentsByParentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAssignedDeploymentsByParent ...
func (c GoldenGateConnectionsClient) ListAssignedDeploymentsByParent(ctx context.Context, id GoldenGateConnectionId) (result ListAssignedDeploymentsByParentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAssignedDeploymentsByParentCustomPager{},
		Path:       fmt.Sprintf("%s/assignedDeployments", id.ID()),
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
		Values *[]AssignedDeployment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAssignedDeploymentsByParentComplete retrieves all the results into a single object
func (c GoldenGateConnectionsClient) ListAssignedDeploymentsByParentComplete(ctx context.Context, id GoldenGateConnectionId) (ListAssignedDeploymentsByParentCompleteResult, error) {
	return c.ListAssignedDeploymentsByParentCompleteMatchingPredicate(ctx, id, AssignedDeploymentOperationPredicate{})
}

// ListAssignedDeploymentsByParentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GoldenGateConnectionsClient) ListAssignedDeploymentsByParentCompleteMatchingPredicate(ctx context.Context, id GoldenGateConnectionId, predicate AssignedDeploymentOperationPredicate) (result ListAssignedDeploymentsByParentCompleteResult, err error) {
	items := make([]AssignedDeployment, 0)

	resp, err := c.ListAssignedDeploymentsByParent(ctx, id)
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

	result = ListAssignedDeploymentsByParentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
