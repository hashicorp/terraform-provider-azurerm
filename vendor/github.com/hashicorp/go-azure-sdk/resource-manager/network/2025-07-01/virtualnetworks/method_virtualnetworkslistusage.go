package virtualnetworks

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

type VirtualNetworksListUsageOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualNetworkUsage
}

type VirtualNetworksListUsageCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualNetworkUsage
}

type VirtualNetworksListUsageCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VirtualNetworksListUsageCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VirtualNetworksListUsage ...
func (c VirtualNetworksClient) VirtualNetworksListUsage(ctx context.Context, id commonids.VirtualNetworkId) (result VirtualNetworksListUsageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VirtualNetworksListUsageCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]VirtualNetworkUsage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualNetworksListUsageComplete retrieves all the results into a single object
func (c VirtualNetworksClient) VirtualNetworksListUsageComplete(ctx context.Context, id commonids.VirtualNetworkId) (VirtualNetworksListUsageCompleteResult, error) {
	return c.VirtualNetworksListUsageCompleteMatchingPredicate(ctx, id, VirtualNetworkUsageOperationPredicate{})
}

// VirtualNetworksListUsageCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworksClient) VirtualNetworksListUsageCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualNetworkId, predicate VirtualNetworkUsageOperationPredicate) (result VirtualNetworksListUsageCompleteResult, err error) {
	items := make([]VirtualNetworkUsage, 0)

	resp, err := c.VirtualNetworksListUsage(ctx, id)
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

	result = VirtualNetworksListUsageCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
