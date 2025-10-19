package deploymentstacksatresourcegroup

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

type DeploymentStacksListAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeploymentStack
}

type DeploymentStacksListAtResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeploymentStack
}

type DeploymentStacksListAtResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DeploymentStacksListAtResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DeploymentStacksListAtResourceGroup ...
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksListAtResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result DeploymentStacksListAtResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DeploymentStacksListAtResourceGroupCustomPager{},
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

// DeploymentStacksListAtResourceGroupComplete retrieves all the results into a single object
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksListAtResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (DeploymentStacksListAtResourceGroupCompleteResult, error) {
	return c.DeploymentStacksListAtResourceGroupCompleteMatchingPredicate(ctx, id, DeploymentStackOperationPredicate{})
}

// DeploymentStacksListAtResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksListAtResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DeploymentStackOperationPredicate) (result DeploymentStacksListAtResourceGroupCompleteResult, err error) {
	items := make([]DeploymentStack, 0)

	resp, err := c.DeploymentStacksListAtResourceGroup(ctx, id)
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

	result = DeploymentStacksListAtResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
