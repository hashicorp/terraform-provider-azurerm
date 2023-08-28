package webtestsapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestsListByComponentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WebTest
}

type WebTestsListByComponentCompleteResult struct {
	Items []WebTest
}

// WebTestsListByComponent ...
func (c WebTestsAPIsClient) WebTestsListByComponent(ctx context.Context, id ComponentId) (result WebTestsListByComponentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/webTests", id.ID()),
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

// WebTestsListByComponentComplete retrieves all the results into a single object
func (c WebTestsAPIsClient) WebTestsListByComponentComplete(ctx context.Context, id ComponentId) (WebTestsListByComponentCompleteResult, error) {
	return c.WebTestsListByComponentCompleteMatchingPredicate(ctx, id, WebTestOperationPredicate{})
}

// WebTestsListByComponentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebTestsAPIsClient) WebTestsListByComponentCompleteMatchingPredicate(ctx context.Context, id ComponentId, predicate WebTestOperationPredicate) (result WebTestsListByComponentCompleteResult, err error) {
	items := make([]WebTest, 0)

	resp, err := c.WebTestsListByComponent(ctx, id)
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

	result = WebTestsListByComponentCompleteResult{
		Items: items,
	}
	return
}
