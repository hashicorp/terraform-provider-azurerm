package graphservicesprods

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

type AccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AccountResource
}

type AccountsListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AccountResource
}

type AccountsListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountsListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountsListByResourceGroup ...
func (c GraphservicesprodsClient) AccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result AccountsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AccountsListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.GraphServices/accounts", id.ID()),
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
		Values *[]AccountResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccountsListByResourceGroupComplete retrieves all the results into a single object
func (c GraphservicesprodsClient) AccountsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (AccountsListByResourceGroupCompleteResult, error) {
	return c.AccountsListByResourceGroupCompleteMatchingPredicate(ctx, id, AccountResourceOperationPredicate{})
}

// AccountsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GraphservicesprodsClient) AccountsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AccountResourceOperationPredicate) (result AccountsListByResourceGroupCompleteResult, err error) {
	items := make([]AccountResource, 0)

	resp, err := c.AccountsListByResourceGroup(ctx, id)
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

	result = AccountsListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
