package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableResourcesListTablesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TableGetResults
}

type TableResourcesListTablesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TableGetResults
}

type TableResourcesListTablesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TableResourcesListTablesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TableResourcesListTables ...
func (c OpenapisClient) TableResourcesListTables(ctx context.Context, id DatabaseAccountId) (result TableResourcesListTablesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &TableResourcesListTablesCustomPager{},
		Path:       fmt.Sprintf("%s/tables", id.ID()),
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
		Values *[]TableGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TableResourcesListTablesComplete retrieves all the results into a single object
func (c OpenapisClient) TableResourcesListTablesComplete(ctx context.Context, id DatabaseAccountId) (TableResourcesListTablesCompleteResult, error) {
	return c.TableResourcesListTablesCompleteMatchingPredicate(ctx, id, TableGetResultsOperationPredicate{})
}

// TableResourcesListTablesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) TableResourcesListTablesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate TableGetResultsOperationPredicate) (result TableResourcesListTablesCompleteResult, err error) {
	items := make([]TableGetResults, 0)

	resp, err := c.TableResourcesListTables(ctx, id)
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

	result = TableResourcesListTablesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
