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

type ListProcessThreadsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProcessThreadInfo
}

type ListProcessThreadsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProcessThreadInfo
}

type ListProcessThreadsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListProcessThreadsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListProcessThreadsSlot ...
func (c WebAppsClient) ListProcessThreadsSlot(ctx context.Context, id SlotProcessId) (result ListProcessThreadsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListProcessThreadsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/threads", id.ID()),
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
		Values *[]ProcessThreadInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListProcessThreadsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListProcessThreadsSlotComplete(ctx context.Context, id SlotProcessId) (ListProcessThreadsSlotCompleteResult, error) {
	return c.ListProcessThreadsSlotCompleteMatchingPredicate(ctx, id, ProcessThreadInfoOperationPredicate{})
}

// ListProcessThreadsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListProcessThreadsSlotCompleteMatchingPredicate(ctx context.Context, id SlotProcessId, predicate ProcessThreadInfoOperationPredicate) (result ListProcessThreadsSlotCompleteResult, err error) {
	items := make([]ProcessThreadInfo, 0)

	resp, err := c.ListProcessThreadsSlot(ctx, id)
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

	result = ListProcessThreadsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
