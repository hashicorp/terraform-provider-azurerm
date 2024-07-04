package privatelinkscopesapis

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

type PrivateLinkScopesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AzureMonitorPrivateLinkScope
}

type PrivateLinkScopesListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AzureMonitorPrivateLinkScope
}

type PrivateLinkScopesListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PrivateLinkScopesListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PrivateLinkScopesListByResourceGroup ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result PrivateLinkScopesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PrivateLinkScopesListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/privateLinkScopes", id.ID()),
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
		Values *[]AzureMonitorPrivateLinkScope `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PrivateLinkScopesListByResourceGroupComplete retrieves all the results into a single object
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (PrivateLinkScopesListByResourceGroupCompleteResult, error) {
	return c.PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate(ctx, id, AzureMonitorPrivateLinkScopeOperationPredicate{})
}

// PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AzureMonitorPrivateLinkScopeOperationPredicate) (result PrivateLinkScopesListByResourceGroupCompleteResult, err error) {
	items := make([]AzureMonitorPrivateLinkScope, 0)

	resp, err := c.PrivateLinkScopesListByResourceGroup(ctx, id)
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

	result = PrivateLinkScopesListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
