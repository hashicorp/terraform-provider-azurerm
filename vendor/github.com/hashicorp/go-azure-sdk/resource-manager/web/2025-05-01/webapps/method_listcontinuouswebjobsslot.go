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

type ListContinuousWebJobsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContinuousWebJob
}

type ListContinuousWebJobsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ContinuousWebJob
}

type ListContinuousWebJobsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListContinuousWebJobsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListContinuousWebJobsSlot ...
func (c WebAppsClient) ListContinuousWebJobsSlot(ctx context.Context, id SlotId) (result ListContinuousWebJobsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListContinuousWebJobsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/continuousWebJobs", id.ID()),
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
		Values *[]ContinuousWebJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListContinuousWebJobsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListContinuousWebJobsSlotComplete(ctx context.Context, id SlotId) (ListContinuousWebJobsSlotCompleteResult, error) {
	return c.ListContinuousWebJobsSlotCompleteMatchingPredicate(ctx, id, ContinuousWebJobOperationPredicate{})
}

// ListContinuousWebJobsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListContinuousWebJobsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate ContinuousWebJobOperationPredicate) (result ListContinuousWebJobsSlotCompleteResult, err error) {
	items := make([]ContinuousWebJob, 0)

	resp, err := c.ListContinuousWebJobsSlot(ctx, id)
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

	result = ListContinuousWebJobsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
