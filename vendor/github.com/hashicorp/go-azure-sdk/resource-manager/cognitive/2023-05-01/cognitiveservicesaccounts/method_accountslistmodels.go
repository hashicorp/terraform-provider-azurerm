package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsListModelsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AccountModel
}

type AccountsListModelsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AccountModel
}

type AccountsListModelsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountsListModelsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountsListModels ...
func (c CognitiveServicesAccountsClient) AccountsListModels(ctx context.Context, id AccountId) (result AccountsListModelsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AccountsListModelsCustomPager{},
		Path:       fmt.Sprintf("%s/models", id.ID()),
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
		Values *[]AccountModel `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccountsListModelsComplete retrieves all the results into a single object
func (c CognitiveServicesAccountsClient) AccountsListModelsComplete(ctx context.Context, id AccountId) (AccountsListModelsCompleteResult, error) {
	return c.AccountsListModelsCompleteMatchingPredicate(ctx, id, AccountModelOperationPredicate{})
}

// AccountsListModelsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CognitiveServicesAccountsClient) AccountsListModelsCompleteMatchingPredicate(ctx context.Context, id AccountId, predicate AccountModelOperationPredicate) (result AccountsListModelsCompleteResult, err error) {
	items := make([]AccountModel, 0)

	resp, err := c.AccountsListModels(ctx, id)
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

	result = AccountsListModelsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
