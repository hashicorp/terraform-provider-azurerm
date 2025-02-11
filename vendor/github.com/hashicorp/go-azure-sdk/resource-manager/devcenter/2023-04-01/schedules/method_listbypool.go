package schedules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByPoolOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Schedule
}

type ListByPoolCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Schedule
}

type ListByPoolCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByPoolCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByPool ...
func (c SchedulesClient) ListByPool(ctx context.Context, id PoolId) (result ListByPoolOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByPoolCustomPager{},
		Path:       fmt.Sprintf("%s/schedules", id.ID()),
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
		Values *[]Schedule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByPoolComplete retrieves all the results into a single object
func (c SchedulesClient) ListByPoolComplete(ctx context.Context, id PoolId) (ListByPoolCompleteResult, error) {
	return c.ListByPoolCompleteMatchingPredicate(ctx, id, ScheduleOperationPredicate{})
}

// ListByPoolCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SchedulesClient) ListByPoolCompleteMatchingPredicate(ctx context.Context, id PoolId, predicate ScheduleOperationPredicate) (result ListByPoolCompleteResult, err error) {
	items := make([]Schedule, 0)

	resp, err := c.ListByPool(ctx, id)
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

	result = ListByPoolCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
