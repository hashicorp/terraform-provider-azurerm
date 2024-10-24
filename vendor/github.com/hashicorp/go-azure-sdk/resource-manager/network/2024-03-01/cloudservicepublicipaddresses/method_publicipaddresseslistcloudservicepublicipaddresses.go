package cloudservicepublicipaddresses

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPAddressesListCloudServicePublicIPAddressesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicIPAddress
}

type PublicIPAddressesListCloudServicePublicIPAddressesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PublicIPAddress
}

type PublicIPAddressesListCloudServicePublicIPAddressesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PublicIPAddressesListCloudServicePublicIPAddressesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PublicIPAddressesListCloudServicePublicIPAddresses ...
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServicePublicIPAddresses(ctx context.Context, id ProviderCloudServiceId) (result PublicIPAddressesListCloudServicePublicIPAddressesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PublicIPAddressesListCloudServicePublicIPAddressesCustomPager{},
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

// PublicIPAddressesListCloudServicePublicIPAddressesComplete retrieves all the results into a single object
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServicePublicIPAddressesComplete(ctx context.Context, id ProviderCloudServiceId) (PublicIPAddressesListCloudServicePublicIPAddressesCompleteResult, error) {
	return c.PublicIPAddressesListCloudServicePublicIPAddressesCompleteMatchingPredicate(ctx, id, PublicIPAddressOperationPredicate{})
}

// PublicIPAddressesListCloudServicePublicIPAddressesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CloudServicePublicIPAddressesClient) PublicIPAddressesListCloudServicePublicIPAddressesCompleteMatchingPredicate(ctx context.Context, id ProviderCloudServiceId, predicate PublicIPAddressOperationPredicate) (result PublicIPAddressesListCloudServicePublicIPAddressesCompleteResult, err error) {
	items := make([]PublicIPAddress, 0)

	resp, err := c.PublicIPAddressesListCloudServicePublicIPAddresses(ctx, id)
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

	result = PublicIPAddressesListCloudServicePublicIPAddressesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
