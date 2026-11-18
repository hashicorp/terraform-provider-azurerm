package edgeactionexecutionfilters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByEdgeActionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EdgeActionExecutionFilter
}

type ListByEdgeActionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EdgeActionExecutionFilter
}

type ListByEdgeActionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByEdgeActionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByEdgeAction ...
func (c EdgeActionExecutionFiltersClient) ListByEdgeAction(ctx context.Context, id EdgeActionId) (result ListByEdgeActionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByEdgeActionCustomPager{},
		Path:       fmt.Sprintf("%s/executionFilters", id.ID()),
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
		Values *[]EdgeActionExecutionFilter `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByEdgeActionComplete retrieves all the results into a single object
func (c EdgeActionExecutionFiltersClient) ListByEdgeActionComplete(ctx context.Context, id EdgeActionId) (ListByEdgeActionCompleteResult, error) {
	return c.ListByEdgeActionCompleteMatchingPredicate(ctx, id, EdgeActionExecutionFilterOperationPredicate{})
}

// ListByEdgeActionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EdgeActionExecutionFiltersClient) ListByEdgeActionCompleteMatchingPredicate(ctx context.Context, id EdgeActionId, predicate EdgeActionExecutionFilterOperationPredicate) (result ListByEdgeActionCompleteResult, err error) {
	items := make([]EdgeActionExecutionFilter, 0)

	resp, err := c.ListByEdgeAction(ctx, id)
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

	result = ListByEdgeActionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
