package appserviceenvironments

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

type ListAppServicePlansOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AppServicePlan
}

type ListAppServicePlansCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AppServicePlan
}

type ListAppServicePlansCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAppServicePlansCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAppServicePlans ...
func (c AppServiceEnvironmentsClient) ListAppServicePlans(ctx context.Context, id commonids.AppServiceEnvironmentId) (result ListAppServicePlansOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAppServicePlansCustomPager{},
		Path:       fmt.Sprintf("%s/serverFarms", id.ID()),
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
		Values *[]AppServicePlan `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAppServicePlansComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListAppServicePlansComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (ListAppServicePlansCompleteResult, error) {
	return c.ListAppServicePlansCompleteMatchingPredicate(ctx, id, AppServicePlanOperationPredicate{})
}

// ListAppServicePlansCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListAppServicePlansCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate AppServicePlanOperationPredicate) (result ListAppServicePlansCompleteResult, err error) {
	items := make([]AppServicePlan, 0)

	resp, err := c.ListAppServicePlans(ctx, id)
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

	result = ListAppServicePlansCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
