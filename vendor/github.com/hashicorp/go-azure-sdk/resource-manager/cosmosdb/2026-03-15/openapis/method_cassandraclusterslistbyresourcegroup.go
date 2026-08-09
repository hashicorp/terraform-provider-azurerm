package openapis

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

type CassandraClustersListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterResource
}

type CassandraClustersListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterResource
}

type CassandraClustersListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CassandraClustersListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CassandraClustersListByResourceGroup ...
func (c OpenapisClient) CassandraClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result CassandraClustersListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CassandraClustersListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DocumentDB/cassandraClusters", id.ID()),
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
		Values *[]ClusterResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CassandraClustersListByResourceGroupComplete retrieves all the results into a single object
func (c OpenapisClient) CassandraClustersListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (CassandraClustersListByResourceGroupCompleteResult, error) {
	return c.CassandraClustersListByResourceGroupCompleteMatchingPredicate(ctx, id, ClusterResourceOperationPredicate{})
}

// CassandraClustersListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CassandraClustersListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ClusterResourceOperationPredicate) (result CassandraClustersListByResourceGroupCompleteResult, err error) {
	items := make([]ClusterResource, 0)

	resp, err := c.CassandraClustersListByResourceGroup(ctx, id)
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

	result = CassandraClustersListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
