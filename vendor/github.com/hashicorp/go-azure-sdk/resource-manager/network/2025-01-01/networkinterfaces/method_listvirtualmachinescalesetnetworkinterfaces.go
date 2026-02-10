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

type ListVirtualMachineScaleSetNetworkInterfacesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkInterface
}

type ListVirtualMachineScaleSetNetworkInterfacesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkInterface
}

type ListVirtualMachineScaleSetNetworkInterfacesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListVirtualMachineScaleSetNetworkInterfacesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListVirtualMachineScaleSetNetworkInterfaces ...
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetNetworkInterfaces(ctx context.Context, id VirtualMachineScaleSetId) (result ListVirtualMachineScaleSetNetworkInterfacesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListVirtualMachineScaleSetNetworkInterfacesCustomPager{},
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

// ListVirtualMachineScaleSetNetworkInterfacesComplete retrieves all the results into a single object
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetNetworkInterfacesComplete(ctx context.Context, id VirtualMachineScaleSetId) (ListVirtualMachineScaleSetNetworkInterfacesCompleteResult, error) {
	return c.ListVirtualMachineScaleSetNetworkInterfacesCompleteMatchingPredicate(ctx, id, NetworkInterfaceOperationPredicate{})
}

// ListVirtualMachineScaleSetNetworkInterfacesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkInterfacesClient) ListVirtualMachineScaleSetNetworkInterfacesCompleteMatchingPredicate(ctx context.Context, id VirtualMachineScaleSetId, predicate NetworkInterfaceOperationPredicate) (result ListVirtualMachineScaleSetNetworkInterfacesCompleteResult, err error) {
	items := make([]NetworkInterface, 0)

	resp, err := c.ListVirtualMachineScaleSetNetworkInterfaces(ctx, id)
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

	result = ListVirtualMachineScaleSetNetworkInterfacesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
