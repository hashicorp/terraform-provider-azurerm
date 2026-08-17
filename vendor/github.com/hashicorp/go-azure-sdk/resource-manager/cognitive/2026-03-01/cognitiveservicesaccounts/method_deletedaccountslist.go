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

type DeletedAccountsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Account
}

type DeletedAccountsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Account
}

type DeletedAccountsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DeletedAccountsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DeletedAccountsList ...
func (c CognitiveServicesAccountsClient) DeletedAccountsList(ctx context.Context, id commonids.SubscriptionId) (result DeletedAccountsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DeletedAccountsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/deletedAccounts", id.ID()),
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

// DeletedAccountsListComplete retrieves all the results into a single object
func (c CognitiveServicesAccountsClient) DeletedAccountsListComplete(ctx context.Context, id commonids.SubscriptionId) (DeletedAccountsListCompleteResult, error) {
	return c.DeletedAccountsListCompleteMatchingPredicate(ctx, id, AccountOperationPredicate{})
}

// DeletedAccountsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CognitiveServicesAccountsClient) DeletedAccountsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AccountOperationPredicate) (result DeletedAccountsListCompleteResult, err error) {
	items := make([]Account, 0)

	resp, err := c.DeletedAccountsList(ctx, id)
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

	result = DeletedAccountsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
