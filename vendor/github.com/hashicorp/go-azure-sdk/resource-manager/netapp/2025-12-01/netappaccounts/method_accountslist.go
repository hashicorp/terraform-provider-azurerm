package netappaccounts

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

type AccountsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetAppAccount
}

type AccountsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetAppAccount
}

type AccountsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountsList ...
func (c NetAppAccountsClient) AccountsList(ctx context.Context, id commonids.ResourceGroupId) (result AccountsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AccountsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.NetApp/netAppAccounts", id.ID()),
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
		Values *[]NetAppAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccountsListComplete retrieves all the results into a single object
func (c NetAppAccountsClient) AccountsListComplete(ctx context.Context, id commonids.ResourceGroupId) (AccountsListCompleteResult, error) {
	return c.AccountsListCompleteMatchingPredicate(ctx, id, NetAppAccountOperationPredicate{})
}

// AccountsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetAppAccountsClient) AccountsListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate NetAppAccountOperationPredicate) (result AccountsListCompleteResult, err error) {
	items := make([]NetAppAccount, 0)

	resp, err := c.AccountsList(ctx, id)
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

	result = AccountsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
