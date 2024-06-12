package hdinsights

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

type ClusterPoolsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterPool
}

type ClusterPoolsListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterPool
}

// ClusterPoolsListBySubscription ...
func (c HdinsightsClient) ClusterPoolsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result ClusterPoolsListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.HDInsight/clusterPools", id.ID()),
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
		Values *[]ClusterPool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClusterPoolsListBySubscriptionComplete retrieves all the results into a single object
func (c HdinsightsClient) ClusterPoolsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ClusterPoolsListBySubscriptionCompleteResult, error) {
	return c.ClusterPoolsListBySubscriptionCompleteMatchingPredicate(ctx, id, ClusterPoolOperationPredicate{})
}

// ClusterPoolsListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) ClusterPoolsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ClusterPoolOperationPredicate) (result ClusterPoolsListBySubscriptionCompleteResult, err error) {
	items := make([]ClusterPool, 0)

	resp, err := c.ClusterPoolsListBySubscription(ctx, id)
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

	result = ClusterPoolsListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
