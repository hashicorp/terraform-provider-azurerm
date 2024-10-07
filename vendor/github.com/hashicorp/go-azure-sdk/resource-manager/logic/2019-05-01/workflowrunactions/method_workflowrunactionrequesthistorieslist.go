package workflowrunactions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionRequestHistoriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RequestHistory
}

type WorkflowRunActionRequestHistoriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RequestHistory
}

type WorkflowRunActionRequestHistoriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkflowRunActionRequestHistoriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkflowRunActionRequestHistoriesList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesList(ctx context.Context, id ActionId) (result WorkflowRunActionRequestHistoriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WorkflowRunActionRequestHistoriesListCustomPager{},
		Path:       fmt.Sprintf("%s/requestHistories", id.ID()),
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
		Values *[]RequestHistory `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkflowRunActionRequestHistoriesListComplete retrieves all the results into a single object
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesListComplete(ctx context.Context, id ActionId) (WorkflowRunActionRequestHistoriesListCompleteResult, error) {
	return c.WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate(ctx, id, RequestHistoryOperationPredicate{})
}

// WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate(ctx context.Context, id ActionId, predicate RequestHistoryOperationPredicate) (result WorkflowRunActionRequestHistoriesListCompleteResult, err error) {
	items := make([]RequestHistory, 0)

	resp, err := c.WorkflowRunActionRequestHistoriesList(ctx, id)
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

	result = WorkflowRunActionRequestHistoriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
