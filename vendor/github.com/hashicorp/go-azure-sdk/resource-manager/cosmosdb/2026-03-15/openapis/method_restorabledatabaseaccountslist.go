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

type RestorableDatabaseAccountsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableDatabaseAccountGetResult
}

type RestorableDatabaseAccountsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableDatabaseAccountGetResult
}

type RestorableDatabaseAccountsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableDatabaseAccountsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableDatabaseAccountsList ...
func (c OpenapisClient) RestorableDatabaseAccountsList(ctx context.Context, id commonids.SubscriptionId) (result RestorableDatabaseAccountsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RestorableDatabaseAccountsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DocumentDB/restorableDatabaseAccounts", id.ID()),
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

// RestorableDatabaseAccountsListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableDatabaseAccountsListComplete(ctx context.Context, id commonids.SubscriptionId) (RestorableDatabaseAccountsListCompleteResult, error) {
	return c.RestorableDatabaseAccountsListCompleteMatchingPredicate(ctx, id, RestorableDatabaseAccountGetResultOperationPredicate{})
}

// RestorableDatabaseAccountsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableDatabaseAccountsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate RestorableDatabaseAccountGetResultOperationPredicate) (result RestorableDatabaseAccountsListCompleteResult, err error) {
	items := make([]RestorableDatabaseAccountGetResult, 0)

	resp, err := c.RestorableDatabaseAccountsList(ctx, id)
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

	result = RestorableDatabaseAccountsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
