package webapps

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

type ListProductionSiteDeploymentStatusesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmDeploymentStatus
}

type ListProductionSiteDeploymentStatusesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmDeploymentStatus
}

type ListProductionSiteDeploymentStatusesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListProductionSiteDeploymentStatusesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListProductionSiteDeploymentStatuses ...
func (c WebAppsClient) ListProductionSiteDeploymentStatuses(ctx context.Context, id commonids.AppServiceId) (result ListProductionSiteDeploymentStatusesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListProductionSiteDeploymentStatusesCustomPager{},
		Path:       fmt.Sprintf("%s/deploymentStatus", id.ID()),
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
		Values *[]CsmDeploymentStatus `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListProductionSiteDeploymentStatusesComplete retrieves all the results into a single object
func (c WebAppsClient) ListProductionSiteDeploymentStatusesComplete(ctx context.Context, id commonids.AppServiceId) (ListProductionSiteDeploymentStatusesCompleteResult, error) {
	return c.ListProductionSiteDeploymentStatusesCompleteMatchingPredicate(ctx, id, CsmDeploymentStatusOperationPredicate{})
}

// ListProductionSiteDeploymentStatusesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListProductionSiteDeploymentStatusesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate CsmDeploymentStatusOperationPredicate) (result ListProductionSiteDeploymentStatusesCompleteResult, err error) {
	items := make([]CsmDeploymentStatus, 0)

	resp, err := c.ListProductionSiteDeploymentStatuses(ctx, id)
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

	result = ListProductionSiteDeploymentStatusesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
