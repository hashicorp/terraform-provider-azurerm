package fleets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetspaceListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FleetspaceResource
}

type FleetspaceListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FleetspaceResource
}

type FleetspaceListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FleetspaceListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FleetspaceList ...
func (c FleetsClient) FleetspaceList(ctx context.Context, id FleetId) (result FleetspaceListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FleetspaceListCustomPager{},
		Path:       fmt.Sprintf("%s/fleetspaces", id.ID()),
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
		Values *[]FleetspaceResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FleetspaceListComplete retrieves all the results into a single object
func (c FleetsClient) FleetspaceListComplete(ctx context.Context, id FleetId) (FleetspaceListCompleteResult, error) {
	return c.FleetspaceListCompleteMatchingPredicate(ctx, id, FleetspaceResourceOperationPredicate{})
}

// FleetspaceListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FleetsClient) FleetspaceListCompleteMatchingPredicate(ctx context.Context, id FleetId, predicate FleetspaceResourceOperationPredicate) (result FleetspaceListCompleteResult, err error) {
	items := make([]FleetspaceResource, 0)

	resp, err := c.FleetspaceList(ctx, id)
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

	result = FleetspaceListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
