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

type FleetListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FleetResource
}

type FleetListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FleetResource
}

type FleetListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FleetListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FleetList ...
func (c FleetsClient) FleetList(ctx context.Context, id commonids.SubscriptionId) (result FleetListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FleetListCustomPager{},
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

// FleetListComplete retrieves all the results into a single object
func (c FleetsClient) FleetListComplete(ctx context.Context, id commonids.SubscriptionId) (FleetListCompleteResult, error) {
	return c.FleetListCompleteMatchingPredicate(ctx, id, FleetResourceOperationPredicate{})
}

// FleetListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FleetsClient) FleetListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate FleetResourceOperationPredicate) (result FleetListCompleteResult, err error) {
	items := make([]FleetResource, 0)

	resp, err := c.FleetList(ctx, id)
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

	result = FleetListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
