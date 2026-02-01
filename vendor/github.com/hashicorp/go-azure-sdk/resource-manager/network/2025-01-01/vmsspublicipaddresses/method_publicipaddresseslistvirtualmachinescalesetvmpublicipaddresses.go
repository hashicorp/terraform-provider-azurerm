package vmsspublicipaddresses

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

type PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicIPAddress
}

type PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PublicIPAddress
}

type PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddresses ...
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddresses(ctx context.Context, id commonids.VirtualMachineScaleSetIPConfigurationId) (result PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCustomPager{},
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

// PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesComplete retrieves all the results into a single object
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesComplete(ctx context.Context, id commonids.VirtualMachineScaleSetIPConfigurationId) (PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteResult, error) {
	return c.PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteMatchingPredicate(ctx, id, PublicIPAddressOperationPredicate{})
}

// PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualMachineScaleSetIPConfigurationId, predicate PublicIPAddressOperationPredicate) (result PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteResult, err error) {
	items := make([]PublicIPAddress, 0)

	resp, err := c.PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddresses(ctx, id)
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

	result = PublicIPAddressesListVirtualMachineScaleSetVMPublicIPAddressesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
