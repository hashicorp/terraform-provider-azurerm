package patchschedules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByRedisResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RedisPatchSchedule
}

type ListByRedisResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RedisPatchSchedule
}

type ListByRedisResourceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByRedisResourceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByRedisResource ...
func (c PatchSchedulesClient) ListByRedisResource(ctx context.Context, id RediId) (result ListByRedisResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByRedisResourceCustomPager{},
		Path:       fmt.Sprintf("%s/patchSchedules", id.ID()),
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
		Values *[]RedisPatchSchedule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByRedisResourceComplete retrieves all the results into a single object
func (c PatchSchedulesClient) ListByRedisResourceComplete(ctx context.Context, id RediId) (ListByRedisResourceCompleteResult, error) {
	return c.ListByRedisResourceCompleteMatchingPredicate(ctx, id, RedisPatchScheduleOperationPredicate{})
}

// ListByRedisResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PatchSchedulesClient) ListByRedisResourceCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisPatchScheduleOperationPredicate) (result ListByRedisResourceCompleteResult, err error) {
	items := make([]RedisPatchSchedule, 0)

	resp, err := c.ListByRedisResource(ctx, id)
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

	result = ListByRedisResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
