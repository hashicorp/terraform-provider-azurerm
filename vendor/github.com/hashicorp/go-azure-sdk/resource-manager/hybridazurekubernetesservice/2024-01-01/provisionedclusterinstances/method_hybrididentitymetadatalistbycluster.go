package provisionedclusterinstances

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

type HybridIdentityMetadataListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HybridIdentityMetadata
}

type HybridIdentityMetadataListByClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HybridIdentityMetadata
}

type HybridIdentityMetadataListByClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *HybridIdentityMetadataListByClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// HybridIdentityMetadataListByCluster ...
func (c ProvisionedClusterInstancesClient) HybridIdentityMetadataListByCluster(ctx context.Context, id commonids.ScopeId) (result HybridIdentityMetadataListByClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &HybridIdentityMetadataListByClusterCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default/hybridIdentityMetadata", id.ID()),
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
		Values *[]HybridIdentityMetadata `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// HybridIdentityMetadataListByClusterComplete retrieves all the results into a single object
func (c ProvisionedClusterInstancesClient) HybridIdentityMetadataListByClusterComplete(ctx context.Context, id commonids.ScopeId) (HybridIdentityMetadataListByClusterCompleteResult, error) {
	return c.HybridIdentityMetadataListByClusterCompleteMatchingPredicate(ctx, id, HybridIdentityMetadataOperationPredicate{})
}

// HybridIdentityMetadataListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvisionedClusterInstancesClient) HybridIdentityMetadataListByClusterCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate HybridIdentityMetadataOperationPredicate) (result HybridIdentityMetadataListByClusterCompleteResult, err error) {
	items := make([]HybridIdentityMetadata, 0)

	resp, err := c.HybridIdentityMetadataListByCluster(ctx, id)
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

	result = HybridIdentityMetadataListByClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
