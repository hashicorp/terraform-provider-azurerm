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

type WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RequestHistory
}

type WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RequestHistory
}

type WorkflowRunActionRepetitionsRequestHistoriesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkflowRunActionRepetitionsRequestHistoriesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkflowRunActionRepetitionsRequestHistoriesList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesList(ctx context.Context, id RepetitionId) (result WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WorkflowRunActionRepetitionsRequestHistoriesListCustomPager{},
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

// WorkflowRunActionRepetitionsRequestHistoriesListComplete retrieves all the results into a single object
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesListComplete(ctx context.Context, id RepetitionId) (WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult, error) {
	return c.WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate(ctx, id, RequestHistoryOperationPredicate{})
}

// WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate(ctx context.Context, id RepetitionId, predicate RequestHistoryOperationPredicate) (result WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult, err error) {
	items := make([]RequestHistory, 0)

	resp, err := c.WorkflowRunActionRepetitionsRequestHistoriesList(ctx, id)
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

	result = WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
