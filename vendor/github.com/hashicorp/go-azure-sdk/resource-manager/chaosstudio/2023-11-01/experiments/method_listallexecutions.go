package experiments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAllExecutionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ExperimentExecution
}

type ListAllExecutionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ExperimentExecution
}

type ListAllExecutionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAllExecutionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAllExecutions ...
func (c ExperimentsClient) ListAllExecutions(ctx context.Context, id ExperimentId) (result ListAllExecutionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAllExecutionsCustomPager{},
		Path:       fmt.Sprintf("%s/executions", id.ID()),
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
		Values *[]ExperimentExecution `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAllExecutionsComplete retrieves all the results into a single object
func (c ExperimentsClient) ListAllExecutionsComplete(ctx context.Context, id ExperimentId) (ListAllExecutionsCompleteResult, error) {
	return c.ListAllExecutionsCompleteMatchingPredicate(ctx, id, ExperimentExecutionOperationPredicate{})
}

// ListAllExecutionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExperimentsClient) ListAllExecutionsCompleteMatchingPredicate(ctx context.Context, id ExperimentId, predicate ExperimentExecutionOperationPredicate) (result ListAllExecutionsCompleteResult, err error) {
	items := make([]ExperimentExecution, 0)

	resp, err := c.ListAllExecutions(ctx, id)
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

	result = ListAllExecutionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
