package hdinsights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolUpgradeHistoriesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterPoolUpgradeHistory
}

type ClusterPoolUpgradeHistoriesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterPoolUpgradeHistory
}

// ClusterPoolUpgradeHistoriesList ...
func (c HdinsightsClient) ClusterPoolUpgradeHistoriesList(ctx context.Context, id ClusterPoolId) (result ClusterPoolUpgradeHistoriesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/upgradeHistories", id.ID()),
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
		Values *[]ClusterPoolUpgradeHistory `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClusterPoolUpgradeHistoriesListComplete retrieves all the results into a single object
func (c HdinsightsClient) ClusterPoolUpgradeHistoriesListComplete(ctx context.Context, id ClusterPoolId) (ClusterPoolUpgradeHistoriesListCompleteResult, error) {
	return c.ClusterPoolUpgradeHistoriesListCompleteMatchingPredicate(ctx, id, ClusterPoolUpgradeHistoryOperationPredicate{})
}

// ClusterPoolUpgradeHistoriesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) ClusterPoolUpgradeHistoriesListCompleteMatchingPredicate(ctx context.Context, id ClusterPoolId, predicate ClusterPoolUpgradeHistoryOperationPredicate) (result ClusterPoolUpgradeHistoriesListCompleteResult, err error) {
	items := make([]ClusterPoolUpgradeHistory, 0)

	resp, err := c.ClusterPoolUpgradeHistoriesList(ctx, id)
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

	result = ClusterPoolUpgradeHistoriesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
