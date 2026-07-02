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

type ListProcessModulesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessModuleInfo
}

type ListProcessModulesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessModuleInfo
}

type ListProcessModulesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListProcessModulesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListProcessModulesSlot ...
func (c WebAppsClient) ListProcessModulesSlot(ctx context.Context, id SlotProcessId) (result ListProcessModulesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListProcessModulesSlotCustomPager{},
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

// ListProcessModulesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListProcessModulesSlotComplete(ctx context.Context, id SlotProcessId) (ListProcessModulesSlotCompleteResult, error) {
	return c.ListProcessModulesSlotCompleteMatchingPredicate(ctx, id, ProcessModuleInfoOperationPredicate{})
}

// ListProcessModulesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListProcessModulesSlotCompleteMatchingPredicate(ctx context.Context, id SlotProcessId, predicate ProcessModuleInfoOperationPredicate) (result ListProcessModulesSlotCompleteResult, err error) {
	items := make([]ProcessModuleInfo, 0)

	resp, err := c.ListProcessModulesSlot(ctx, id)
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

	result = ListProcessModulesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
