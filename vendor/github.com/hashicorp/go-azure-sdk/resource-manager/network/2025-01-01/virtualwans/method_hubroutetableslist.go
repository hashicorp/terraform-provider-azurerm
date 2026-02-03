package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubRouteTablesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HubRouteTable
}

type HubRouteTablesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HubRouteTable
}

type HubRouteTablesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *HubRouteTablesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// HubRouteTablesList ...
func (c VirtualWANsClient) HubRouteTablesList(ctx context.Context, id VirtualHubId) (result HubRouteTablesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &HubRouteTablesListCustomPager{},
		Path:       fmt.Sprintf("%s/hubRouteTables", id.ID()),
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
		Values *[]HubRouteTable `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// HubRouteTablesListComplete retrieves all the results into a single object
func (c VirtualWANsClient) HubRouteTablesListComplete(ctx context.Context, id VirtualHubId) (HubRouteTablesListCompleteResult, error) {
	return c.HubRouteTablesListCompleteMatchingPredicate(ctx, id, HubRouteTableOperationPredicate{})
}

// HubRouteTablesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) HubRouteTablesListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate HubRouteTableOperationPredicate) (result HubRouteTablesListCompleteResult, err error) {
	items := make([]HubRouteTable, 0)

	resp, err := c.HubRouteTablesList(ctx, id)
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

	result = HubRouteTablesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
