package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSlotSiteDeploymentStatusesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmDeploymentStatus
}

type ListSlotSiteDeploymentStatusesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmDeploymentStatus
}

type ListSlotSiteDeploymentStatusesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSlotSiteDeploymentStatusesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSlotSiteDeploymentStatusesSlot ...
func (c WebAppsClient) ListSlotSiteDeploymentStatusesSlot(ctx context.Context, id SlotId) (result ListSlotSiteDeploymentStatusesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSlotSiteDeploymentStatusesSlotCustomPager{},
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

// ListSlotSiteDeploymentStatusesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSlotSiteDeploymentStatusesSlotComplete(ctx context.Context, id SlotId) (ListSlotSiteDeploymentStatusesSlotCompleteResult, error) {
	return c.ListSlotSiteDeploymentStatusesSlotCompleteMatchingPredicate(ctx, id, CsmDeploymentStatusOperationPredicate{})
}

// ListSlotSiteDeploymentStatusesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSlotSiteDeploymentStatusesSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate CsmDeploymentStatusOperationPredicate) (result ListSlotSiteDeploymentStatusesSlotCompleteResult, err error) {
	items := make([]CsmDeploymentStatus, 0)

	resp, err := c.ListSlotSiteDeploymentStatusesSlot(ctx, id)
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

	result = ListSlotSiteDeploymentStatusesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
