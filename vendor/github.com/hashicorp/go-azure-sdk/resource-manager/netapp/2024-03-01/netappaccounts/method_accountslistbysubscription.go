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

type AccountsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetAppAccount
}

type AccountsListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetAppAccount
}

type AccountsListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountsListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountsListBySubscription ...
func (c NetAppAccountsClient) AccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result AccountsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AccountsListBySubscriptionCustomPager{},
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

// AccountsListBySubscriptionComplete retrieves all the results into a single object
func (c NetAppAccountsClient) AccountsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (AccountsListBySubscriptionCompleteResult, error) {
	return c.AccountsListBySubscriptionCompleteMatchingPredicate(ctx, id, NetAppAccountOperationPredicate{})
}

// AccountsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetAppAccountsClient) AccountsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate NetAppAccountOperationPredicate) (result AccountsListBySubscriptionCompleteResult, err error) {
	items := make([]NetAppAccount, 0)

	resp, err := c.AccountsListBySubscription(ctx, id)
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

	result = AccountsListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
