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

type ListTriggeredWebJobsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TriggeredWebJob
}

type ListTriggeredWebJobsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TriggeredWebJob
}

type ListTriggeredWebJobsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListTriggeredWebJobsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListTriggeredWebJobsSlot ...
func (c WebAppsClient) ListTriggeredWebJobsSlot(ctx context.Context, id SlotId) (result ListTriggeredWebJobsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListTriggeredWebJobsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/triggeredWebJobs", id.ID()),
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
		Values *[]TriggeredWebJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListTriggeredWebJobsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListTriggeredWebJobsSlotComplete(ctx context.Context, id SlotId) (ListTriggeredWebJobsSlotCompleteResult, error) {
	return c.ListTriggeredWebJobsSlotCompleteMatchingPredicate(ctx, id, TriggeredWebJobOperationPredicate{})
}

// ListTriggeredWebJobsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListTriggeredWebJobsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate TriggeredWebJobOperationPredicate) (result ListTriggeredWebJobsSlotCompleteResult, err error) {
	items := make([]TriggeredWebJob, 0)

	resp, err := c.ListTriggeredWebJobsSlot(ctx, id)
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

	result = ListTriggeredWebJobsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
