package firewallresources

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

type FirewallsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FirewallResource
}

type FirewallsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FirewallResource
}

type FirewallsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FirewallsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FirewallsListByResourceGroup ...
func (c FirewallResourcesClient) FirewallsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result FirewallsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FirewallsListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/PaloAltoNetworks.Cloudngfw/firewalls", id.ID()),
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
		Values *[]FirewallResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FirewallsListByResourceGroupComplete retrieves all the results into a single object
func (c FirewallResourcesClient) FirewallsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (FirewallsListByResourceGroupCompleteResult, error) {
	return c.FirewallsListByResourceGroupCompleteMatchingPredicate(ctx, id, FirewallResourceOperationPredicate{})
}

// FirewallsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FirewallResourcesClient) FirewallsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate FirewallResourceOperationPredicate) (result FirewallsListByResourceGroupCompleteResult, err error) {
	items := make([]FirewallResource, 0)

	resp, err := c.FirewallsListByResourceGroup(ctx, id)
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

	result = FirewallsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
