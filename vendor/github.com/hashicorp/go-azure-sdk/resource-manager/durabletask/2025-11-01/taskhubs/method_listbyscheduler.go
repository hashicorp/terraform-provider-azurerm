package taskhubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySchedulerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TaskHub
}

type ListBySchedulerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TaskHub
}

type ListBySchedulerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySchedulerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByScheduler ...
func (c TaskHubsClient) ListByScheduler(ctx context.Context, id SchedulerId) (result ListBySchedulerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBySchedulerCustomPager{},
		Path:       fmt.Sprintf("%s/taskHubs", id.ID()),
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
		Values *[]TaskHub `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySchedulerComplete retrieves all the results into a single object
func (c TaskHubsClient) ListBySchedulerComplete(ctx context.Context, id SchedulerId) (ListBySchedulerCompleteResult, error) {
	return c.ListBySchedulerCompleteMatchingPredicate(ctx, id, TaskHubOperationPredicate{})
}

// ListBySchedulerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TaskHubsClient) ListBySchedulerCompleteMatchingPredicate(ctx context.Context, id SchedulerId, predicate TaskHubOperationPredicate) (result ListBySchedulerCompleteResult, err error) {
	items := make([]TaskHub, 0)

	resp, err := c.ListByScheduler(ctx, id)
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

	result = ListBySchedulerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
