package resourceproviders

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

type ListAseRegionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AseRegion
}

type ListAseRegionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AseRegion
}

type ListAseRegionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAseRegionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAseRegions ...
func (c ResourceProvidersClient) ListAseRegions(ctx context.Context, id commonids.SubscriptionId) (result ListAseRegionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAseRegionsCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Web/aseRegions", id.ID()),
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
		Values *[]AseRegion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAseRegionsComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListAseRegionsComplete(ctx context.Context, id commonids.SubscriptionId) (ListAseRegionsCompleteResult, error) {
	return c.ListAseRegionsCompleteMatchingPredicate(ctx, id, AseRegionOperationPredicate{})
}

// ListAseRegionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListAseRegionsCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AseRegionOperationPredicate) (result ListAseRegionsCompleteResult, err error) {
	items := make([]AseRegion, 0)

	resp, err := c.ListAseRegions(ctx, id)
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

	result = ListAseRegionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
