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

type ListCapacitiesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StampCapacity
}

type ListCapacitiesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StampCapacity
}

type ListCapacitiesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCapacitiesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListCapacities ...
func (c AppServiceEnvironmentsClient) ListCapacities(ctx context.Context, id commonids.AppServiceEnvironmentId) (result ListCapacitiesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListCapacitiesCustomPager{},
		Path:       fmt.Sprintf("%s/capacities/compute", id.ID()),
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
		Values *[]StampCapacity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListCapacitiesComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListCapacitiesComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (ListCapacitiesCompleteResult, error) {
	return c.ListCapacitiesCompleteMatchingPredicate(ctx, id, StampCapacityOperationPredicate{})
}

// ListCapacitiesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListCapacitiesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate StampCapacityOperationPredicate) (result ListCapacitiesCompleteResult, err error) {
	items := make([]StampCapacity, 0)

	resp, err := c.ListCapacities(ctx, id)
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

	result = ListCapacitiesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
