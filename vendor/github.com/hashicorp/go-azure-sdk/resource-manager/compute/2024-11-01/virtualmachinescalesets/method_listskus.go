package virtualmachinescalesets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualMachineScaleSetSku
}

type ListSkusCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VirtualMachineScaleSetSku
}

type ListSkusCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSkusCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSkus ...
func (c VirtualMachineScaleSetsClient) ListSkus(ctx context.Context, id VirtualMachineScaleSetId) (result ListSkusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSkusCustomPager{},
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]VirtualMachineScaleSetSku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSkusComplete retrieves all the results into a single object
func (c VirtualMachineScaleSetsClient) ListSkusComplete(ctx context.Context, id VirtualMachineScaleSetId) (ListSkusCompleteResult, error) {
	return c.ListSkusCompleteMatchingPredicate(ctx, id, VirtualMachineScaleSetSkuOperationPredicate{})
}

// ListSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualMachineScaleSetsClient) ListSkusCompleteMatchingPredicate(ctx context.Context, id VirtualMachineScaleSetId, predicate VirtualMachineScaleSetSkuOperationPredicate) (result ListSkusCompleteResult, err error) {
	items := make([]VirtualMachineScaleSetSku, 0)

	resp, err := c.ListSkus(ctx, id)
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

	result = ListSkusCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
