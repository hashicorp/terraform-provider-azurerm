package eventhubsclusters

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

type ClustersListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Cluster
}

type ClustersListBySubscriptionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Cluster
}

type ClustersListBySubscriptionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ClustersListBySubscriptionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ClustersListBySubscription ...
func (c EventHubsClustersClient) ClustersListBySubscription(ctx context.Context, id commonids.SubscriptionId) (result ClustersListBySubscriptionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ClustersListBySubscriptionCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.EventHub/clusters", id.ID()),
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
		Values *[]Cluster `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ClustersListBySubscriptionComplete retrieves all the results into a single object
func (c EventHubsClustersClient) ClustersListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ClustersListBySubscriptionCompleteResult, error) {
	return c.ClustersListBySubscriptionCompleteMatchingPredicate(ctx, id, ClusterOperationPredicate{})
}

// ClustersListBySubscriptionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventHubsClustersClient) ClustersListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ClusterOperationPredicate) (result ClustersListBySubscriptionCompleteResult, err error) {
	items := make([]Cluster, 0)

	resp, err := c.ClustersListBySubscription(ctx, id)
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

	result = ClustersListBySubscriptionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
