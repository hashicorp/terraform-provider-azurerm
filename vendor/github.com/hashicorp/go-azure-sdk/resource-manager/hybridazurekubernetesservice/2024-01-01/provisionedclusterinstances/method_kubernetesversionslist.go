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

type KubernetesVersionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]KubernetesVersionProfile
}

type KubernetesVersionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []KubernetesVersionProfile
}

type KubernetesVersionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *KubernetesVersionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// KubernetesVersionsList ...
func (c ProvisionedClusterInstancesClient) KubernetesVersionsList(ctx context.Context, id commonids.ScopeId) (result KubernetesVersionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &KubernetesVersionsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.HybridContainerService/kubernetesVersions", id.ID()),
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
		Values *[]KubernetesVersionProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// KubernetesVersionsListComplete retrieves all the results into a single object
func (c ProvisionedClusterInstancesClient) KubernetesVersionsListComplete(ctx context.Context, id commonids.ScopeId) (KubernetesVersionsListCompleteResult, error) {
	return c.KubernetesVersionsListCompleteMatchingPredicate(ctx, id, KubernetesVersionProfileOperationPredicate{})
}

// KubernetesVersionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvisionedClusterInstancesClient) KubernetesVersionsListCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate KubernetesVersionProfileOperationPredicate) (result KubernetesVersionsListCompleteResult, err error) {
	items := make([]KubernetesVersionProfile, 0)

	resp, err := c.KubernetesVersionsList(ctx, id)
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

	result = KubernetesVersionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
