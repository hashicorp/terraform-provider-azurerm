package networkinterfaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListVirtualMachineScaleSetVMNetworkInterfacesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkInterface
}

type ListVirtualMachineScaleSetVMNetworkInterfacesCompleteResult struct {
	Items []NetworkInterface
}

// ListVirtualMachineScaleSetVMNetworkInterfaces ...
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetVMNetworkInterfaces(ctx context.Context, id VirtualMachineId) (result ListVirtualMachineScaleSetVMNetworkInterfacesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/networkInterfaces", id.ID()),
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
		Values *[]NetworkInterface `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListVirtualMachineScaleSetVMNetworkInterfacesComplete retrieves all the results into a single object
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetVMNetworkInterfacesComplete(ctx context.Context, id VirtualMachineId) (ListVirtualMachineScaleSetVMNetworkInterfacesCompleteResult, error) {
	return c.ListVirtualMachineScaleSetVMNetworkInterfacesCompleteMatchingPredicate(ctx, id, NetworkInterfaceOperationPredicate{})
}

// ListVirtualMachineScaleSetVMNetworkInterfacesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetVMNetworkInterfacesCompleteMatchingPredicate(ctx context.Context, id VirtualMachineId, predicate NetworkInterfaceOperationPredicate) (result ListVirtualMachineScaleSetVMNetworkInterfacesCompleteResult, err error) {
	items := make([]NetworkInterface, 0)

	resp, err := c.ListVirtualMachineScaleSetVMNetworkInterfaces(ctx, id)
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

	result = ListVirtualMachineScaleSetVMNetworkInterfacesCompleteResult{
		Items: items,
	}
	return
}
