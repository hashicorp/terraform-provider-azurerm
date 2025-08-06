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

type WorkflowRunActionRepetitionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkflowRunActionRepetitionDefinition
}

type WorkflowRunActionRepetitionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WorkflowRunActionRepetitionDefinition
}

type WorkflowRunActionRepetitionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkflowRunActionRepetitionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkflowRunActionRepetitionsList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsList(ctx context.Context, id ActionId) (result WorkflowRunActionRepetitionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WorkflowRunActionRepetitionsListCustomPager{},
		Path:       fmt.Sprintf("%s/repetitions", id.ID()),
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
		Values *[]WorkflowRunActionRepetitionDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkflowRunActionRepetitionsListComplete retrieves all the results into a single object
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsListComplete(ctx context.Context, id ActionId) (WorkflowRunActionRepetitionsListCompleteResult, error) {
	return c.WorkflowRunActionRepetitionsListCompleteMatchingPredicate(ctx, id, WorkflowRunActionRepetitionDefinitionOperationPredicate{})
}

// WorkflowRunActionRepetitionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsListCompleteMatchingPredicate(ctx context.Context, id ActionId, predicate WorkflowRunActionRepetitionDefinitionOperationPredicate) (result WorkflowRunActionRepetitionsListCompleteResult, err error) {
	items := make([]WorkflowRunActionRepetitionDefinition, 0)

	resp, err := c.WorkflowRunActionRepetitionsList(ctx, id)
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

	result = WorkflowRunActionRepetitionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
