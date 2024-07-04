package resourceguards

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

type GetResourcesInResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceGuardResource
}

type GetResourcesInResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceGuardResource
}

type GetResourcesInResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetResourcesInResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetResourcesInResourceGroup ...
func (c ResourceGuardsClient) GetResourcesInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result GetResourcesInResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetResourcesInResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DataProtection/resourceGuards", id.ID()),
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
		Values *[]ResourceGuardResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetResourcesInResourceGroupComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetResourcesInResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GetResourcesInResourceGroupCompleteResult, error) {
	return c.GetResourcesInResourceGroupCompleteMatchingPredicate(ctx, id, ResourceGuardResourceOperationPredicate{})
}

// GetResourcesInResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetResourcesInResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ResourceGuardResourceOperationPredicate) (result GetResourcesInResourceGroupCompleteResult, err error) {
	items := make([]ResourceGuardResource, 0)

	resp, err := c.GetResourcesInResourceGroup(ctx, id)
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

	result = GetResourcesInResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
