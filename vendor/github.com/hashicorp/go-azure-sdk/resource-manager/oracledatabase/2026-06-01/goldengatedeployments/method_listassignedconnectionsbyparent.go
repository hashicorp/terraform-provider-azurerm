package goldengatedeployments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAssignedConnectionsByParentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AssignedConnection
}

type ListAssignedConnectionsByParentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AssignedConnection
}

type ListAssignedConnectionsByParentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAssignedConnectionsByParentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAssignedConnectionsByParent ...
func (c GoldenGateDeploymentsClient) ListAssignedConnectionsByParent(ctx context.Context, id GoldenGateDeploymentId) (result ListAssignedConnectionsByParentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAssignedConnectionsByParentCustomPager{},
		Path:       fmt.Sprintf("%s/assignedConnections", id.ID()),
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
		Values *[]AssignedConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAssignedConnectionsByParentComplete retrieves all the results into a single object
func (c GoldenGateDeploymentsClient) ListAssignedConnectionsByParentComplete(ctx context.Context, id GoldenGateDeploymentId) (ListAssignedConnectionsByParentCompleteResult, error) {
	return c.ListAssignedConnectionsByParentCompleteMatchingPredicate(ctx, id, AssignedConnectionOperationPredicate{})
}

// ListAssignedConnectionsByParentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GoldenGateDeploymentsClient) ListAssignedConnectionsByParentCompleteMatchingPredicate(ctx context.Context, id GoldenGateDeploymentId, predicate AssignedConnectionOperationPredicate) (result ListAssignedConnectionsByParentCompleteResult, err error) {
	items := make([]AssignedConnection, 0)

	resp, err := c.ListAssignedConnectionsByParent(ctx, id)
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

	result = ListAssignedConnectionsByParentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
