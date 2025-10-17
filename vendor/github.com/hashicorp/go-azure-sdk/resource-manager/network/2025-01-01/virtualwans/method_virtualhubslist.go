package virtualwans

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

type VirtualHubsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualHub
}

type VirtualHubsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualHub
}

type VirtualHubsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VirtualHubsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VirtualHubsList ...
func (c VirtualWANsClient) VirtualHubsList(ctx context.Context, id commonids.SubscriptionId) (result VirtualHubsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VirtualHubsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/virtualHubs", id.ID()),
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
		Values *[]VirtualHub `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualHubsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualHubsListComplete(ctx context.Context, id commonids.SubscriptionId) (VirtualHubsListCompleteResult, error) {
	return c.VirtualHubsListCompleteMatchingPredicate(ctx, id, VirtualHubOperationPredicate{})
}

// VirtualHubsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualHubsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate VirtualHubOperationPredicate) (result VirtualHubsListCompleteResult, err error) {
	items := make([]VirtualHub, 0)

	resp, err := c.VirtualHubsList(ctx, id)
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

	result = VirtualHubsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
