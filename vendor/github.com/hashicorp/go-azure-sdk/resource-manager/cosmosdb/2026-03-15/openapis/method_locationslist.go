package openapis

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

type LocationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LocationGetResult
}

type LocationsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LocationGetResult
}

type LocationsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocationsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocationsList ...
func (c OpenapisClient) LocationsList(ctx context.Context, id commonids.SubscriptionId) (result LocationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LocationsListCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.DocumentDB/locations", id.ID()),
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
		Values *[]LocationGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocationsListComplete retrieves all the results into a single object
func (c OpenapisClient) LocationsListComplete(ctx context.Context, id commonids.SubscriptionId) (LocationsListCompleteResult, error) {
	return c.LocationsListCompleteMatchingPredicate(ctx, id, LocationGetResultOperationPredicate{})
}

// LocationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) LocationsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate LocationGetResultOperationPredicate) (result LocationsListCompleteResult, err error) {
	items := make([]LocationGetResult, 0)

	resp, err := c.LocationsList(ctx, id)
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

	result = LocationsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
