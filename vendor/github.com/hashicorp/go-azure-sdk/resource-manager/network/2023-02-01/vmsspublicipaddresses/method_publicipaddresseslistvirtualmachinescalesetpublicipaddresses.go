package vmsspublicipaddresses

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicIPAddress
}

type PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteResult struct {
	Items []PublicIPAddress
}

// PublicIPAddressesListVirtualMachineScaleSetPublicIPAddresses ...
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetPublicIPAddresses(ctx context.Context, id VirtualMachineScaleSetId) (result PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesOperationResponse, err error) {
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

// PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesComplete retrieves all the results into a single object
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesComplete(ctx context.Context, id VirtualMachineScaleSetId) (PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteResult, error) {
	return c.PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteMatchingPredicate(ctx, id, PublicIPAddressOperationPredicate{})
}

// PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VMSSPublicIPAddressesClient) PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteMatchingPredicate(ctx context.Context, id VirtualMachineScaleSetId, predicate PublicIPAddressOperationPredicate) (result PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteResult, err error) {
	items := make([]PublicIPAddress, 0)

	resp, err := c.PublicIPAddressesListVirtualMachineScaleSetPublicIPAddresses(ctx, id)
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

	result = PublicIPAddressesListVirtualMachineScaleSetPublicIPAddressesCompleteResult{
		Items: items,
	}
	return
}
