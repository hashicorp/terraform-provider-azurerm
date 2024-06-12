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

type ClusterAvailableUpgradesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterAvailableUpgrade
}

type ClusterAvailableUpgradesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterAvailableUpgrade
}

// ClusterAvailableUpgradesList ...
func (c HdinsightsClient) ClusterAvailableUpgradesList(ctx context.Context, id ClusterId) (result ClusterAvailableUpgradesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/availableUpgrades", id.ID()),
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
		Values *[]ClusterAvailableUpgrade `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClusterAvailableUpgradesListComplete retrieves all the results into a single object
func (c HdinsightsClient) ClusterAvailableUpgradesListComplete(ctx context.Context, id ClusterId) (ClusterAvailableUpgradesListCompleteResult, error) {
	return c.ClusterAvailableUpgradesListCompleteMatchingPredicate(ctx, id, ClusterAvailableUpgradeOperationPredicate{})
}

// ClusterAvailableUpgradesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) ClusterAvailableUpgradesListCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate ClusterAvailableUpgradeOperationPredicate) (result ClusterAvailableUpgradesListCompleteResult, err error) {
	items := make([]ClusterAvailableUpgrade, 0)

	resp, err := c.ClusterAvailableUpgradesList(ctx, id)
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

	result = ClusterAvailableUpgradesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
