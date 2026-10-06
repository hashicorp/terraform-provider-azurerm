package staticsites

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

type GetStaticSitesByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteARMResource
}

type GetStaticSitesByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteARMResource
}

type GetStaticSitesByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetStaticSitesByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetStaticSitesByResourceGroup ...
func (c StaticSitesClient) GetStaticSitesByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result GetStaticSitesByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetStaticSitesByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Web/staticSites", id.ID()),
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
		Values *[]StaticSiteARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetStaticSitesByResourceGroupComplete retrieves all the results into a single object
func (c StaticSitesClient) GetStaticSitesByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GetStaticSitesByResourceGroupCompleteResult, error) {
	return c.GetStaticSitesByResourceGroupCompleteMatchingPredicate(ctx, id, StaticSiteARMResourceOperationPredicate{})
}

// GetStaticSitesByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetStaticSitesByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate StaticSiteARMResourceOperationPredicate) (result GetStaticSitesByResourceGroupCompleteResult, err error) {
	items := make([]StaticSiteARMResource, 0)

	resp, err := c.GetStaticSitesByResourceGroup(ctx, id)
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

	result = GetStaticSitesByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
