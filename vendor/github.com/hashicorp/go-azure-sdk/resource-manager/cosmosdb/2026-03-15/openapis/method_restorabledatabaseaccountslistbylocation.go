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

type RestorableDatabaseAccountsListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableDatabaseAccountGetResult
}

type RestorableDatabaseAccountsListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableDatabaseAccountGetResult
}

type RestorableDatabaseAccountsListByLocationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableDatabaseAccountsListByLocationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableDatabaseAccountsListByLocation ...
func (c OpenapisClient) RestorableDatabaseAccountsListByLocation(ctx context.Context, id LocationId) (result RestorableDatabaseAccountsListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RestorableDatabaseAccountsListByLocationCustomPager{},
		Path:       fmt.Sprintf("%s/restorableDatabaseAccounts", id.ID()),
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
		Values *[]RestorableDatabaseAccountGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableDatabaseAccountsListByLocationComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableDatabaseAccountsListByLocationComplete(ctx context.Context, id LocationId) (RestorableDatabaseAccountsListByLocationCompleteResult, error) {
	return c.RestorableDatabaseAccountsListByLocationCompleteMatchingPredicate(ctx, id, RestorableDatabaseAccountGetResultOperationPredicate{})
}

// RestorableDatabaseAccountsListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableDatabaseAccountsListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate RestorableDatabaseAccountGetResultOperationPredicate) (result RestorableDatabaseAccountsListByLocationCompleteResult, err error) {
	items := make([]RestorableDatabaseAccountGetResult, 0)

	resp, err := c.RestorableDatabaseAccountsListByLocation(ctx, id)
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

	result = RestorableDatabaseAccountsListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
