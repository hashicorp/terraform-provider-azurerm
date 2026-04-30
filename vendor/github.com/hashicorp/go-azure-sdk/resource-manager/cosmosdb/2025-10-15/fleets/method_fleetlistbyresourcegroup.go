package fleets

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

type FleetListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FleetResource
}

type FleetListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FleetResource
}

type FleetListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FleetListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FleetListByResourceGroup ...
func (c FleetsClient) FleetListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result FleetListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FleetListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DocumentDB/fleets", id.ID()),
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
		Values *[]FleetResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FleetListByResourceGroupComplete retrieves all the results into a single object
func (c FleetsClient) FleetListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (FleetListByResourceGroupCompleteResult, error) {
	return c.FleetListByResourceGroupCompleteMatchingPredicate(ctx, id, FleetResourceOperationPredicate{})
}

// FleetListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FleetsClient) FleetListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate FleetResourceOperationPredicate) (result FleetListByResourceGroupCompleteResult, err error) {
	items := make([]FleetResource, 0)

	resp, err := c.FleetListByResourceGroup(ctx, id)
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

	result = FleetListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
