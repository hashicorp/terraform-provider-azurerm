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

type FleetspaceAccountListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FleetspaceAccountResource
}

type FleetspaceAccountListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FleetspaceAccountResource
}

type FleetspaceAccountListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FleetspaceAccountListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FleetspaceAccountList ...
func (c FleetsClient) FleetspaceAccountList(ctx context.Context, id FleetspaceId) (result FleetspaceAccountListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FleetspaceAccountListCustomPager{},
		Path:       fmt.Sprintf("%s/fleetspaceAccounts", id.ID()),
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
		Values *[]FleetspaceAccountResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FleetspaceAccountListComplete retrieves all the results into a single object
func (c FleetsClient) FleetspaceAccountListComplete(ctx context.Context, id FleetspaceId) (FleetspaceAccountListCompleteResult, error) {
	return c.FleetspaceAccountListCompleteMatchingPredicate(ctx, id, FleetspaceAccountResourceOperationPredicate{})
}

// FleetspaceAccountListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FleetsClient) FleetspaceAccountListCompleteMatchingPredicate(ctx context.Context, id FleetspaceId, predicate FleetspaceAccountResourceOperationPredicate) (result FleetspaceAccountListCompleteResult, err error) {
	items := make([]FleetspaceAccountResource, 0)

	resp, err := c.FleetspaceAccountList(ctx, id)
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

	result = FleetspaceAccountListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
