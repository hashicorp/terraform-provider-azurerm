package maintenanceconfigurations

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

type ListByManagedClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MaintenanceConfiguration
}

type ListByManagedClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MaintenanceConfiguration
}

type ListByManagedClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByManagedClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByManagedCluster ...
func (c MaintenanceConfigurationsClient) ListByManagedCluster(ctx context.Context, id commonids.KubernetesClusterId) (result ListByManagedClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByManagedClusterCustomPager{},
		Path:       fmt.Sprintf("%s/maintenanceConfigurations", id.ID()),
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
		Values *[]MaintenanceConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByManagedClusterComplete retrieves all the results into a single object
func (c MaintenanceConfigurationsClient) ListByManagedClusterComplete(ctx context.Context, id commonids.KubernetesClusterId) (ListByManagedClusterCompleteResult, error) {
	return c.ListByManagedClusterCompleteMatchingPredicate(ctx, id, MaintenanceConfigurationOperationPredicate{})
}

// ListByManagedClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MaintenanceConfigurationsClient) ListByManagedClusterCompleteMatchingPredicate(ctx context.Context, id commonids.KubernetesClusterId, predicate MaintenanceConfigurationOperationPredicate) (result ListByManagedClusterCompleteResult, err error) {
	items := make([]MaintenanceConfiguration, 0)

	resp, err := c.ListByManagedCluster(ctx, id)
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

	result = ListByManagedClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
