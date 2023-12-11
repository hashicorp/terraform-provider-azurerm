package cognitiveservicesaccounts

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
	Model        *[]Account
}

type AccountsListCompleteResult struct {
	Items []Account
}

// AccountsList ...
func (c CognitiveServicesAccountsClient) AccountsList(ctx context.Context, id commonids.SubscriptionId) (result AccountsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/accounts", id.ID()),
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
		Values *[]Account `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccountsListComplete retrieves all the results into a single object
func (c CognitiveServicesAccountsClient) AccountsListComplete(ctx context.Context, id commonids.SubscriptionId) (AccountsListCompleteResult, error) {
	return c.AccountsListCompleteMatchingPredicate(ctx, id, AccountOperationPredicate{})
}

// AccountsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CognitiveServicesAccountsClient) AccountsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AccountOperationPredicate) (result AccountsListCompleteResult, err error) {
	items := make([]Account, 0)

	resp, err := c.AccountsList(ctx, id)
	if err != nil {
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
		Items: items,
	}
	return
}
