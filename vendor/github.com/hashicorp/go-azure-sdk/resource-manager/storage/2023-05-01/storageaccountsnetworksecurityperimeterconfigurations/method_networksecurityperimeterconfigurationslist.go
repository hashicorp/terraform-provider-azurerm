package storageaccountsnetworksecurityperimeterconfigurations

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

type NetworkSecurityPerimeterConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkSecurityPerimeterConfiguration
}

type NetworkSecurityPerimeterConfigurationsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkSecurityPerimeterConfiguration
}

type NetworkSecurityPerimeterConfigurationsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NetworkSecurityPerimeterConfigurationsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NetworkSecurityPerimeterConfigurationsList ...
func (c StorageAccountsNetworkSecurityPerimeterConfigurationsClient) NetworkSecurityPerimeterConfigurationsList(ctx context.Context, id commonids.StorageAccountId) (result NetworkSecurityPerimeterConfigurationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NetworkSecurityPerimeterConfigurationsListCustomPager{},
		Path:       fmt.Sprintf("%s/networkSecurityPerimeterConfigurations", id.ID()),
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
		Values *[]NetworkSecurityPerimeterConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NetworkSecurityPerimeterConfigurationsListComplete retrieves all the results into a single object
func (c StorageAccountsNetworkSecurityPerimeterConfigurationsClient) NetworkSecurityPerimeterConfigurationsListComplete(ctx context.Context, id commonids.StorageAccountId) (NetworkSecurityPerimeterConfigurationsListCompleteResult, error) {
	return c.NetworkSecurityPerimeterConfigurationsListCompleteMatchingPredicate(ctx, id, NetworkSecurityPerimeterConfigurationOperationPredicate{})
}

// NetworkSecurityPerimeterConfigurationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageAccountsNetworkSecurityPerimeterConfigurationsClient) NetworkSecurityPerimeterConfigurationsListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, predicate NetworkSecurityPerimeterConfigurationOperationPredicate) (result NetworkSecurityPerimeterConfigurationsListCompleteResult, err error) {
	items := make([]NetworkSecurityPerimeterConfiguration, 0)

	resp, err := c.NetworkSecurityPerimeterConfigurationsList(ctx, id)
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

	result = NetworkSecurityPerimeterConfigurationsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
