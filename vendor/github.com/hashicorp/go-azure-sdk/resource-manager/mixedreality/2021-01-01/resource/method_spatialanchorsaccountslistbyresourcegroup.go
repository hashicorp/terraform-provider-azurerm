package resource

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

type SpatialAnchorsAccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SpatialAnchorsAccount
}

type SpatialAnchorsAccountsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SpatialAnchorsAccount
}

type SpatialAnchorsAccountsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SpatialAnchorsAccountsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SpatialAnchorsAccountsListByResourceGroup ...
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result SpatialAnchorsAccountsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SpatialAnchorsAccountsListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts", id.ID()),
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
		Values *[]SpatialAnchorsAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SpatialAnchorsAccountsListByResourceGroupComplete retrieves all the results into a single object
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (SpatialAnchorsAccountsListByResourceGroupCompleteResult, error) {
	return c.SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate(ctx, id, SpatialAnchorsAccountOperationPredicate{})
}

// SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate SpatialAnchorsAccountOperationPredicate) (result SpatialAnchorsAccountsListByResourceGroupCompleteResult, err error) {
	items := make([]SpatialAnchorsAccount, 0)

	resp, err := c.SpatialAnchorsAccountsListByResourceGroup(ctx, id)
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

	result = SpatialAnchorsAccountsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
