package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseAccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseAccountGetResults
}

type DatabaseAccountsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatabaseAccountGetResults
}

type DatabaseAccountsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseAccountsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseAccountsListByResourceGroup ...
func (c OpenapisClient) DatabaseAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result DatabaseAccountsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DatabaseAccountsListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DocumentDB/databaseAccounts", id.ID()),
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
		Values *[]DatabaseAccountGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DatabaseAccountsListByResourceGroupComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseAccountsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (DatabaseAccountsListByResourceGroupCompleteResult, error) {
	return c.DatabaseAccountsListByResourceGroupCompleteMatchingPredicate(ctx, id, DatabaseAccountGetResultsOperationPredicate{})
}

// DatabaseAccountsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseAccountsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DatabaseAccountGetResultsOperationPredicate) (result DatabaseAccountsListByResourceGroupCompleteResult, err error) {
	items := make([]DatabaseAccountGetResults, 0)

	resp, err := c.DatabaseAccountsListByResourceGroup(ctx, id)
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

	result = DatabaseAccountsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
