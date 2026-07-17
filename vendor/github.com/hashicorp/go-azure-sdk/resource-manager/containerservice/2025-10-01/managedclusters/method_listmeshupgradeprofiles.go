package managedclusters

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

type ListMeshUpgradeProfilesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MeshUpgradeProfile
}

type ListMeshUpgradeProfilesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MeshUpgradeProfile
}

type ListMeshUpgradeProfilesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListMeshUpgradeProfilesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListMeshUpgradeProfiles ...
func (c ManagedClustersClient) ListMeshUpgradeProfiles(ctx context.Context, id commonids.KubernetesClusterId) (result ListMeshUpgradeProfilesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListMeshUpgradeProfilesCustomPager{},
		Path:       fmt.Sprintf("%s/meshUpgradeProfiles", id.ID()),
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
		Values *[]MeshUpgradeProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListMeshUpgradeProfilesComplete retrieves all the results into a single object
func (c ManagedClustersClient) ListMeshUpgradeProfilesComplete(ctx context.Context, id commonids.KubernetesClusterId) (ListMeshUpgradeProfilesCompleteResult, error) {
	return c.ListMeshUpgradeProfilesCompleteMatchingPredicate(ctx, id, MeshUpgradeProfileOperationPredicate{})
}

// ListMeshUpgradeProfilesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedClustersClient) ListMeshUpgradeProfilesCompleteMatchingPredicate(ctx context.Context, id commonids.KubernetesClusterId, predicate MeshUpgradeProfileOperationPredicate) (result ListMeshUpgradeProfilesCompleteResult, err error) {
	items := make([]MeshUpgradeProfile, 0)

	resp, err := c.ListMeshUpgradeProfiles(ctx, id)
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

	result = ListMeshUpgradeProfilesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
