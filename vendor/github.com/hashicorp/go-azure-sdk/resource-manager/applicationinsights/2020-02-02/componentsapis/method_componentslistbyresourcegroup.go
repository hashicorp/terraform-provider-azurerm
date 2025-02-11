package componentsapis

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

type ComponentsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationInsightsComponent
}

type ComponentsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApplicationInsightsComponent
}

type ComponentsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ComponentsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ComponentsListByResourceGroup ...
func (c ComponentsAPIsClient) ComponentsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result ComponentsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ComponentsListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/components", id.ID()),
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
		Values *[]ApplicationInsightsComponent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ComponentsListByResourceGroupComplete retrieves all the results into a single object
func (c ComponentsAPIsClient) ComponentsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ComponentsListByResourceGroupCompleteResult, error) {
	return c.ComponentsListByResourceGroupCompleteMatchingPredicate(ctx, id, ApplicationInsightsComponentOperationPredicate{})
}

// ComponentsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ComponentsAPIsClient) ComponentsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ApplicationInsightsComponentOperationPredicate) (result ComponentsListByResourceGroupCompleteResult, err error) {
	items := make([]ApplicationInsightsComponent, 0)

	resp, err := c.ComponentsListByResourceGroup(ctx, id)
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

	result = ComponentsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
