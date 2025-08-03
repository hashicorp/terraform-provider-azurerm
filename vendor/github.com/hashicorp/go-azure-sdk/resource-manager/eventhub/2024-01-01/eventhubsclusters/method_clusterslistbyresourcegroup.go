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

type ClustersListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Cluster
}

type ClustersListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Cluster
}

type ClustersListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ClustersListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ClustersListByResourceGroup ...
func (c EventHubsClustersClient) ClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result ClustersListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ClustersListByResourceGroupCustomPager{},
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

// ClustersListByResourceGroupComplete retrieves all the results into a single object
func (c EventHubsClustersClient) ClustersListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ClustersListByResourceGroupCompleteResult, error) {
	return c.ClustersListByResourceGroupCompleteMatchingPredicate(ctx, id, ClusterOperationPredicate{})
}

// ClustersListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EventHubsClustersClient) ClustersListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ClusterOperationPredicate) (result ClustersListByResourceGroupCompleteResult, err error) {
	items := make([]Cluster, 0)

	resp, err := c.ClustersListByResourceGroup(ctx, id)
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

	result = ClustersListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
