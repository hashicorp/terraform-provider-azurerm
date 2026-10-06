package listdeletedrunbooks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountListDeletedRunbooksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedRunbook
}

type AutomationAccountListDeletedRunbooksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedRunbook
}

type AutomationAccountListDeletedRunbooksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AutomationAccountListDeletedRunbooksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AutomationAccountListDeletedRunbooks ...
func (c ListDeletedRunbooksClient) AutomationAccountListDeletedRunbooks(ctx context.Context, id AutomationAccountId) (result AutomationAccountListDeletedRunbooksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &AutomationAccountListDeletedRunbooksCustomPager{},
		Path:       fmt.Sprintf("%s/listDeletedRunbooks", id.ID()),
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
		Values *[]DeletedRunbook `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AutomationAccountListDeletedRunbooksComplete retrieves all the results into a single object
func (c ListDeletedRunbooksClient) AutomationAccountListDeletedRunbooksComplete(ctx context.Context, id AutomationAccountId) (AutomationAccountListDeletedRunbooksCompleteResult, error) {
	return c.AutomationAccountListDeletedRunbooksCompleteMatchingPredicate(ctx, id, DeletedRunbookOperationPredicate{})
}

// AutomationAccountListDeletedRunbooksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ListDeletedRunbooksClient) AutomationAccountListDeletedRunbooksCompleteMatchingPredicate(ctx context.Context, id AutomationAccountId, predicate DeletedRunbookOperationPredicate) (result AutomationAccountListDeletedRunbooksCompleteResult, err error) {
	items := make([]DeletedRunbook, 0)

	resp, err := c.AutomationAccountListDeletedRunbooks(ctx, id)
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

	result = AutomationAccountListDeletedRunbooksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
