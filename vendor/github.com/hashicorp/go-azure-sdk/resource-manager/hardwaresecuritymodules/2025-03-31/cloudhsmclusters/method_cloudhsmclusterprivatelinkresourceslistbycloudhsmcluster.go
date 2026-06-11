package cloudhsmclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateLinkResource
}

type CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PrivateLinkResource
}

type CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CloudHsmClusterPrivateLinkResourcesListByCloudHsmCluster ...
func (c CloudHsmClustersClient) CloudHsmClusterPrivateLinkResourcesListByCloudHsmCluster(ctx context.Context, id CloudHsmClusterId) (result CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCustomPager{},
		Path:       fmt.Sprintf("%s/privateLinkResources", id.ID()),
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
		Values *[]PrivateLinkResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterComplete retrieves all the results into a single object
func (c CloudHsmClustersClient) CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterComplete(ctx context.Context, id CloudHsmClusterId) (CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteResult, error) {
	return c.CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteMatchingPredicate(ctx, id, PrivateLinkResourceOperationPredicate{})
}

// CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CloudHsmClustersClient) CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteMatchingPredicate(ctx context.Context, id CloudHsmClusterId, predicate PrivateLinkResourceOperationPredicate) (result CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteResult, err error) {
	items := make([]PrivateLinkResource, 0)

	resp, err := c.CloudHsmClusterPrivateLinkResourcesListByCloudHsmCluster(ctx, id)
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

	result = CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
