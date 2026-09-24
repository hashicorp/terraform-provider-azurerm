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

type VirtualHubBgpConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BgpConnection
}

type VirtualHubBgpConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BgpConnection
}

type VirtualHubBgpConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VirtualHubBgpConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VirtualHubBgpConnectionsList ...
func (c VirtualWANsClient) VirtualHubBgpConnectionsList(ctx context.Context, id VirtualHubId) (result VirtualHubBgpConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VirtualHubBgpConnectionsListCustomPager{},
		Path:       fmt.Sprintf("%s/bgpConnections", id.ID()),
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
		Values *[]BgpConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualHubBgpConnectionsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualHubBgpConnectionsListComplete(ctx context.Context, id VirtualHubId) (VirtualHubBgpConnectionsListCompleteResult, error) {
	return c.VirtualHubBgpConnectionsListCompleteMatchingPredicate(ctx, id, BgpConnectionOperationPredicate{})
}

// VirtualHubBgpConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualHubBgpConnectionsListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate BgpConnectionOperationPredicate) (result VirtualHubBgpConnectionsListCompleteResult, err error) {
	items := make([]BgpConnection, 0)

	resp, err := c.VirtualHubBgpConnectionsList(ctx, id)
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

	result = VirtualHubBgpConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
