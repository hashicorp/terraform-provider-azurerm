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

type AgentPoolListByProvisionedClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AgentPool
}

type AgentPoolListByProvisionedClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AgentPool
}

type AgentPoolListByProvisionedClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AgentPoolListByProvisionedClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AgentPoolListByProvisionedCluster ...
func (c ProvisionedClusterInstancesClient) AgentPoolListByProvisionedCluster(ctx context.Context, id commonids.ScopeId) (result AgentPoolListByProvisionedClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AgentPoolListByProvisionedClusterCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.HybridContainerService/provisionedClusterInstances/default/agentPools", id.ID()),
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
		Values *[]AgentPool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AgentPoolListByProvisionedClusterComplete retrieves all the results into a single object
func (c ProvisionedClusterInstancesClient) AgentPoolListByProvisionedClusterComplete(ctx context.Context, id commonids.ScopeId) (AgentPoolListByProvisionedClusterCompleteResult, error) {
	return c.AgentPoolListByProvisionedClusterCompleteMatchingPredicate(ctx, id, AgentPoolOperationPredicate{})
}

// AgentPoolListByProvisionedClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvisionedClusterInstancesClient) AgentPoolListByProvisionedClusterCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate AgentPoolOperationPredicate) (result AgentPoolListByProvisionedClusterCompleteResult, err error) {
	items := make([]AgentPool, 0)

	resp, err := c.AgentPoolListByProvisionedCluster(ctx, id)
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

	result = AgentPoolListByProvisionedClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
