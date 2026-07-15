package resourceguardresources

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

type ResourceGuardsGetResourcesInResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceGuardResource
}

type ResourceGuardsGetResourcesInResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceGuardResource
}

type ResourceGuardsGetResourcesInResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ResourceGuardsGetResourcesInResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ResourceGuardsGetResourcesInResourceGroup ...
func (c ResourceGuardResourcesClient) ResourceGuardsGetResourcesInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result ResourceGuardsGetResourcesInResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ResourceGuardsGetResourcesInResourceGroupCustomPager{},
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

// ResourceGuardsGetResourcesInResourceGroupComplete retrieves all the results into a single object
func (c ResourceGuardResourcesClient) ResourceGuardsGetResourcesInResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ResourceGuardsGetResourcesInResourceGroupCompleteResult, error) {
	return c.ResourceGuardsGetResourcesInResourceGroupCompleteMatchingPredicate(ctx, id, ResourceGuardResourceOperationPredicate{})
}

// ResourceGuardsGetResourcesInResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardResourcesClient) ResourceGuardsGetResourcesInResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ResourceGuardResourceOperationPredicate) (result ResourceGuardsGetResourcesInResourceGroupCompleteResult, err error) {
	items := make([]ResourceGuardResource, 0)

	resp, err := c.ResourceGuardsGetResourcesInResourceGroup(ctx, id)
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

	result = ResourceGuardsGetResourcesInResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
