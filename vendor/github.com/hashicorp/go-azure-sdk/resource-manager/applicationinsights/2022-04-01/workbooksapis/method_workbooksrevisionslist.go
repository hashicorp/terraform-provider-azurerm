package workbooksapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksRevisionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Workbook
}

type WorkbooksRevisionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Workbook
}

type WorkbooksRevisionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *WorkbooksRevisionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// WorkbooksRevisionsList ...
func (c WorkbooksAPIsClient) WorkbooksRevisionsList(ctx context.Context, id WorkbookId) (result WorkbooksRevisionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &WorkbooksRevisionsListCustomPager{},
		Path:       fmt.Sprintf("%s/revisions", id.ID()),
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
		Values *[]Workbook `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkbooksRevisionsListComplete retrieves all the results into a single object
func (c WorkbooksAPIsClient) WorkbooksRevisionsListComplete(ctx context.Context, id WorkbookId) (WorkbooksRevisionsListCompleteResult, error) {
	return c.WorkbooksRevisionsListCompleteMatchingPredicate(ctx, id, WorkbookOperationPredicate{})
}

// WorkbooksRevisionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WorkbooksAPIsClient) WorkbooksRevisionsListCompleteMatchingPredicate(ctx context.Context, id WorkbookId, predicate WorkbookOperationPredicate) (result WorkbooksRevisionsListCompleteResult, err error) {
	items := make([]Workbook, 0)

	resp, err := c.WorkbooksRevisionsList(ctx, id)
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

	result = WorkbooksRevisionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
