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

type ClusterPoolAvailableUpgradesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterPoolAvailableUpgrade
}

type ClusterPoolAvailableUpgradesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterPoolAvailableUpgrade
}

// ClusterPoolAvailableUpgradesList ...
func (c HdinsightsClient) ClusterPoolAvailableUpgradesList(ctx context.Context, id ClusterPoolId) (result ClusterPoolAvailableUpgradesListOperationResponse, err error) {
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
		Values *[]ClusterPoolAvailableUpgrade `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClusterPoolAvailableUpgradesListComplete retrieves all the results into a single object
func (c HdinsightsClient) ClusterPoolAvailableUpgradesListComplete(ctx context.Context, id ClusterPoolId) (ClusterPoolAvailableUpgradesListCompleteResult, error) {
	return c.ClusterPoolAvailableUpgradesListCompleteMatchingPredicate(ctx, id, ClusterPoolAvailableUpgradeOperationPredicate{})
}

// ClusterPoolAvailableUpgradesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) ClusterPoolAvailableUpgradesListCompleteMatchingPredicate(ctx context.Context, id ClusterPoolId, predicate ClusterPoolAvailableUpgradeOperationPredicate) (result ClusterPoolAvailableUpgradesListCompleteResult, err error) {
	items := make([]ClusterPoolAvailableUpgrade, 0)

	resp, err := c.ClusterPoolAvailableUpgradesList(ctx, id)
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

	result = ClusterPoolAvailableUpgradesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
