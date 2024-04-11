package connectedclusters

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

type ConnectedClusterListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectedCluster
}

type ConnectedClusterListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectedCluster
}

// ConnectedClusterListBySubscription ...
func (c ConnectedClustersClient) ConnectedClusterListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result ConnectedClusterListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Kubernetes/connectedClusters", id.ID()),
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
		Values *[]ConnectedCluster `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConnectedClusterListBySubscriptionComplete retrieves all the results into a single object
func (c ConnectedClustersClient) ConnectedClusterListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ConnectedClusterListBySubscriptionCompleteResult, error) {
	return c.ConnectedClusterListBySubscriptionCompleteMatchingPredicate(ctx, id, ConnectedClusterOperationPredicate{})
}

// ConnectedClusterListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ConnectedClustersClient) ConnectedClusterListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ConnectedClusterOperationPredicate) (result ConnectedClusterListBySubscriptionCompleteResult, err error) {
	items := make([]ConnectedCluster, 0)

	resp, err := c.ConnectedClusterListBySubscription(ctx, id)
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

	result = ConnectedClusterListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
