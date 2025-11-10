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

type VpnServerConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnServerConfiguration
}

type VpnServerConfigurationsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VpnServerConfiguration
}

type VpnServerConfigurationsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VpnServerConfigurationsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VpnServerConfigurationsList ...
func (c VirtualWANsClient) VpnServerConfigurationsList(ctx context.Context, id commonids.SubscriptionId) (result VpnServerConfigurationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VpnServerConfigurationsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/vpnServerConfigurations", id.ID()),
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
		Values *[]VpnServerConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnServerConfigurationsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnServerConfigurationsListComplete(ctx context.Context, id commonids.SubscriptionId) (VpnServerConfigurationsListCompleteResult, error) {
	return c.VpnServerConfigurationsListCompleteMatchingPredicate(ctx, id, VpnServerConfigurationOperationPredicate{})
}

// VpnServerConfigurationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnServerConfigurationsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate VpnServerConfigurationOperationPredicate) (result VpnServerConfigurationsListCompleteResult, err error) {
	items := make([]VpnServerConfiguration, 0)

	resp, err := c.VpnServerConfigurationsList(ctx, id)
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

	result = VpnServerConfigurationsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
