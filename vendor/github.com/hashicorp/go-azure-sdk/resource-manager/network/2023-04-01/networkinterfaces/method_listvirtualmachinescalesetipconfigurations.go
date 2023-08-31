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

type ListVirtualMachineScaleSetIPConfigurationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkInterfaceIPConfiguration
}

type ListVirtualMachineScaleSetIPConfigurationsCompleteResult struct {
	Items []NetworkInterfaceIPConfiguration
}

type ListVirtualMachineScaleSetIPConfigurationsOperationOptions struct {
	Expand *string
}

func DefaultListVirtualMachineScaleSetIPConfigurationsOperationOptions() ListVirtualMachineScaleSetIPConfigurationsOperationOptions {
	return ListVirtualMachineScaleSetIPConfigurationsOperationOptions{}
}

func (o ListVirtualMachineScaleSetIPConfigurationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListVirtualMachineScaleSetIPConfigurationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListVirtualMachineScaleSetIPConfigurationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

// ListVirtualMachineScaleSetIPConfigurations ...
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetIPConfigurations(ctx context.Context, id commonids.VirtualMachineScaleSetNetworkInterfaceId, options ListVirtualMachineScaleSetIPConfigurationsOperationOptions) (result ListVirtualMachineScaleSetIPConfigurationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/ipConfigurations", id.ID()),
		OptionsObject: options,
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
		Values *[]NetworkInterfaceIPConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListVirtualMachineScaleSetIPConfigurationsComplete retrieves all the results into a single object
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetIPConfigurationsComplete(ctx context.Context, id commonids.VirtualMachineScaleSetNetworkInterfaceId, options ListVirtualMachineScaleSetIPConfigurationsOperationOptions) (ListVirtualMachineScaleSetIPConfigurationsCompleteResult, error) {
	return c.ListVirtualMachineScaleSetIPConfigurationsCompleteMatchingPredicate(ctx, id, options, NetworkInterfaceIPConfigurationOperationPredicate{})
}

// ListVirtualMachineScaleSetIPConfigurationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetIPConfigurationsCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualMachineScaleSetNetworkInterfaceId, options ListVirtualMachineScaleSetIPConfigurationsOperationOptions, predicate NetworkInterfaceIPConfigurationOperationPredicate) (result ListVirtualMachineScaleSetIPConfigurationsCompleteResult, err error) {
	items := make([]NetworkInterfaceIPConfiguration, 0)

	resp, err := c.ListVirtualMachineScaleSetIPConfigurations(ctx, id, options)
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

	result = ListVirtualMachineScaleSetIPConfigurationsCompleteResult{
		Items: items,
	}
	return
}
