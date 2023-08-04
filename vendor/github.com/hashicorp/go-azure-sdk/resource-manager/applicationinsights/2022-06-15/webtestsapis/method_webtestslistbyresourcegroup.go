package webtestsapis

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

type WebTestsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WebTest
}

type WebTestsListByResourceGroupCompleteResult struct {
	Items []WebTest
}

// WebTestsListByResourceGroup ...
func (c WebTestsAPIsClient) WebTestsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result WebTestsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/webTests", id.ID()),
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
		Values *[]WebTest `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WebTestsListByResourceGroupComplete retrieves all the results into a single object
func (c WebTestsAPIsClient) WebTestsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (WebTestsListByResourceGroupCompleteResult, error) {
	return c.WebTestsListByResourceGroupCompleteMatchingPredicate(ctx, id, WebTestOperationPredicate{})
}

// WebTestsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebTestsAPIsClient) WebTestsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate WebTestOperationPredicate) (result WebTestsListByResourceGroupCompleteResult, err error) {
	items := make([]WebTest, 0)

	resp, err := c.WebTestsListByResourceGroup(ctx, id)
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

	result = WebTestsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
