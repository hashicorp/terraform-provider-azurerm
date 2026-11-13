package virtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAvailableSizesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualMachineSize
}

type ListAvailableSizesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualMachineSize
}

type ListAvailableSizesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAvailableSizesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAvailableSizes ...
func (c VirtualMachinesClient) ListAvailableSizes(ctx context.Context, id VirtualMachineId) (result ListAvailableSizesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAvailableSizesCustomPager{},
		Path:       fmt.Sprintf("%s/vmSizes", id.ID()),
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
		Values *[]VirtualMachineSize `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAvailableSizesComplete retrieves all the results into a single object
func (c VirtualMachinesClient) ListAvailableSizesComplete(ctx context.Context, id VirtualMachineId) (ListAvailableSizesCompleteResult, error) {
	return c.ListAvailableSizesCompleteMatchingPredicate(ctx, id, VirtualMachineSizeOperationPredicate{})
}

// ListAvailableSizesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachinesClient) ListAvailableSizesCompleteMatchingPredicate(ctx context.Context, id VirtualMachineId, predicate VirtualMachineSizeOperationPredicate) (result ListAvailableSizesCompleteResult, err error) {
	items := make([]VirtualMachineSize, 0)

	resp, err := c.ListAvailableSizes(ctx, id)
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

	result = ListAvailableSizesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
