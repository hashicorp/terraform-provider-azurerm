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

type HubVirtualNetworkConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HubVirtualNetworkConnection
}

type HubVirtualNetworkConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HubVirtualNetworkConnection
}

type HubVirtualNetworkConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *HubVirtualNetworkConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// HubVirtualNetworkConnectionsList ...
func (c VirtualWANsClient) HubVirtualNetworkConnectionsList(ctx context.Context, id VirtualHubId) (result HubVirtualNetworkConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &HubVirtualNetworkConnectionsListCustomPager{},
		Path:       fmt.Sprintf("%s/hubVirtualNetworkConnections", id.ID()),
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
		Values *[]HubVirtualNetworkConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// HubVirtualNetworkConnectionsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) HubVirtualNetworkConnectionsListComplete(ctx context.Context, id VirtualHubId) (HubVirtualNetworkConnectionsListCompleteResult, error) {
	return c.HubVirtualNetworkConnectionsListCompleteMatchingPredicate(ctx, id, HubVirtualNetworkConnectionOperationPredicate{})
}

// HubVirtualNetworkConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) HubVirtualNetworkConnectionsListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate HubVirtualNetworkConnectionOperationPredicate) (result HubVirtualNetworkConnectionsListCompleteResult, err error) {
	items := make([]HubVirtualNetworkConnection, 0)

	resp, err := c.HubVirtualNetworkConnectionsList(ctx, id)
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

	result = HubVirtualNetworkConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
