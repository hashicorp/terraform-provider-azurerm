package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListInstanceProcessModulesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessModuleInfo
}

type ListInstanceProcessModulesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessModuleInfo
}

type ListInstanceProcessModulesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceProcessModulesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceProcessModulesSlot ...
func (c WebAppsClient) ListInstanceProcessModulesSlot(ctx context.Context, id SlotInstanceProcessId) (result ListInstanceProcessModulesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceProcessModulesSlotCustomPager{},
		Path:       fmt.Sprintf("%s/modules", id.ID()),
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
		Values *[]ProcessModuleInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInstanceProcessModulesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceProcessModulesSlotComplete(ctx context.Context, id SlotInstanceProcessId) (ListInstanceProcessModulesSlotCompleteResult, error) {
	return c.ListInstanceProcessModulesSlotCompleteMatchingPredicate(ctx, id, ProcessModuleInfoOperationPredicate{})
}

// ListInstanceProcessModulesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceProcessModulesSlotCompleteMatchingPredicate(ctx context.Context, id SlotInstanceProcessId, predicate ProcessModuleInfoOperationPredicate) (result ListInstanceProcessModulesSlotCompleteResult, err error) {
	items := make([]ProcessModuleInfo, 0)

	resp, err := c.ListInstanceProcessModulesSlot(ctx, id)
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

	result = ListInstanceProcessModulesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
