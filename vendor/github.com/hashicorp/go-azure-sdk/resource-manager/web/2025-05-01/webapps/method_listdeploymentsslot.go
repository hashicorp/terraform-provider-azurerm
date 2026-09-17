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

type ListDeploymentsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Deployment
}

type ListDeploymentsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Deployment
}

type ListDeploymentsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListDeploymentsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListDeploymentsSlot ...
func (c WebAppsClient) ListDeploymentsSlot(ctx context.Context, id SlotId) (result ListDeploymentsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListDeploymentsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/deployments", id.ID()),
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
		Values *[]Deployment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListDeploymentsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListDeploymentsSlotComplete(ctx context.Context, id SlotId) (ListDeploymentsSlotCompleteResult, error) {
	return c.ListDeploymentsSlotCompleteMatchingPredicate(ctx, id, DeploymentOperationPredicate{})
}

// ListDeploymentsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListDeploymentsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate DeploymentOperationPredicate) (result ListDeploymentsSlotCompleteResult, err error) {
	items := make([]Deployment, 0)

	resp, err := c.ListDeploymentsSlot(ctx, id)
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

	result = ListDeploymentsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
