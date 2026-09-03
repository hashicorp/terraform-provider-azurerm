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

type ListTriggeredWebJobHistoryOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TriggeredJobHistory
}

type ListTriggeredWebJobHistoryCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TriggeredJobHistory
}

type ListTriggeredWebJobHistoryCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListTriggeredWebJobHistoryCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListTriggeredWebJobHistory ...
func (c WebAppsClient) ListTriggeredWebJobHistory(ctx context.Context, id TriggeredWebJobId) (result ListTriggeredWebJobHistoryOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListTriggeredWebJobHistoryCustomPager{},
		Path:       fmt.Sprintf("%s/history", id.ID()),
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
		Values *[]TriggeredJobHistory `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListTriggeredWebJobHistoryComplete retrieves all the results into a single object
func (c WebAppsClient) ListTriggeredWebJobHistoryComplete(ctx context.Context, id TriggeredWebJobId) (ListTriggeredWebJobHistoryCompleteResult, error) {
	return c.ListTriggeredWebJobHistoryCompleteMatchingPredicate(ctx, id, TriggeredJobHistoryOperationPredicate{})
}

// ListTriggeredWebJobHistoryCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListTriggeredWebJobHistoryCompleteMatchingPredicate(ctx context.Context, id TriggeredWebJobId, predicate TriggeredJobHistoryOperationPredicate) (result ListTriggeredWebJobHistoryCompleteResult, err error) {
	items := make([]TriggeredJobHistory, 0)

	resp, err := c.ListTriggeredWebJobHistory(ctx, id)
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

	result = ListTriggeredWebJobHistoryCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
