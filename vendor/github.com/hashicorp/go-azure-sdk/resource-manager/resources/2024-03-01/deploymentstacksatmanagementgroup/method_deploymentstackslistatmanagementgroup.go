package deploymentstacksatmanagementgroup

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

type DeploymentStacksListAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentStack
}

type DeploymentStacksListAtManagementGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentStack
}

type DeploymentStacksListAtManagementGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DeploymentStacksListAtManagementGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DeploymentStacksListAtManagementGroup ...
func (c DeploymentStacksAtManagementGroupClient) DeploymentStacksListAtManagementGroup(ctx context.Context, id commonids.ManagementGroupId) (result DeploymentStacksListAtManagementGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DeploymentStacksListAtManagementGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Resources/deploymentStacks", id.ID()),
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
		Values *[]DeploymentStack `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DeploymentStacksListAtManagementGroupComplete retrieves all the results into a single object
func (c DeploymentStacksAtManagementGroupClient) DeploymentStacksListAtManagementGroupComplete(ctx context.Context, id commonids.ManagementGroupId) (DeploymentStacksListAtManagementGroupCompleteResult, error) {
	return c.DeploymentStacksListAtManagementGroupCompleteMatchingPredicate(ctx, id, DeploymentStackOperationPredicate{})
}

// DeploymentStacksListAtManagementGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentStacksAtManagementGroupClient) DeploymentStacksListAtManagementGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ManagementGroupId, predicate DeploymentStackOperationPredicate) (result DeploymentStacksListAtManagementGroupCompleteResult, err error) {
	items := make([]DeploymentStack, 0)

	resp, err := c.DeploymentStacksListAtManagementGroup(ctx, id)
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

	result = DeploymentStacksListAtManagementGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
