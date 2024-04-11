package cloudservicepublicipaddresses

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

type PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicIPAddress
}

type PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PublicIPAddress
}

// PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses ...
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses(ctx context.Context, id commonids.CloudServicesIPConfigurationId) (result PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/publicIPAddresses", id.ID()),
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
		Values *[]PublicIPAddress `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesComplete retrieves all the results into a single object
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesComplete(ctx context.Context, id commonids.CloudServicesIPConfigurationId) (PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteResult, error) {
	return c.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteMatchingPredicate(ctx, id, PublicIPAddressOperationPredicate{})
}

// PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteMatchingPredicate(ctx context.Context, id commonids.CloudServicesIPConfigurationId, predicate PublicIPAddressOperationPredicate) (result PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteResult, err error) {
	items := make([]PublicIPAddress, 0)

	resp, err := c.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses(ctx, id)
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

	result = PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
