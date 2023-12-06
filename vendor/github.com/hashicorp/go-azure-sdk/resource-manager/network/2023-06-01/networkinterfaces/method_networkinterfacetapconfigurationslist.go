package networkinterfaces

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

type NetworkInterfaceTapConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkInterfaceTapConfiguration
}

type NetworkInterfaceTapConfigurationsListCompleteResult struct {
	Items []NetworkInterfaceTapConfiguration
}

// NetworkInterfaceTapConfigurationsList ...
func (c NetworkInterfacesClient) NetworkInterfaceTapConfigurationsList(ctx context.Context, id commonids.NetworkInterfaceId) (result NetworkInterfaceTapConfigurationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/tapConfigurations", id.ID()),
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
		Values *[]NetworkInterfaceTapConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NetworkInterfaceTapConfigurationsListComplete retrieves all the results into a single object
func (c NetworkInterfacesClient) NetworkInterfaceTapConfigurationsListComplete(ctx context.Context, id commonids.NetworkInterfaceId) (NetworkInterfaceTapConfigurationsListCompleteResult, error) {
	return c.NetworkInterfaceTapConfigurationsListCompleteMatchingPredicate(ctx, id, NetworkInterfaceTapConfigurationOperationPredicate{})
}

// NetworkInterfaceTapConfigurationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkInterfacesClient) NetworkInterfaceTapConfigurationsListCompleteMatchingPredicate(ctx context.Context, id commonids.NetworkInterfaceId, predicate NetworkInterfaceTapConfigurationOperationPredicate) (result NetworkInterfaceTapConfigurationsListCompleteResult, err error) {
	items := make([]NetworkInterfaceTapConfiguration, 0)

	resp, err := c.NetworkInterfaceTapConfigurationsList(ctx, id)
	if err != nil {
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

	result = NetworkInterfaceTapConfigurationsListCompleteResult{
		Items: items,
	}
	return
}
