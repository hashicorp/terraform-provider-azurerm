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

type NetworkVirtualApplianceConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkVirtualApplianceConnection
}

type NetworkVirtualApplianceConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkVirtualApplianceConnection
}

type NetworkVirtualApplianceConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NetworkVirtualApplianceConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NetworkVirtualApplianceConnectionsList ...
func (c VirtualWANsClient) NetworkVirtualApplianceConnectionsList(ctx context.Context, id NetworkVirtualApplianceId) (result NetworkVirtualApplianceConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NetworkVirtualApplianceConnectionsListCustomPager{},
		Path:       fmt.Sprintf("%s/networkVirtualApplianceConnections", id.ID()),
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
		Values *[]NetworkVirtualApplianceConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NetworkVirtualApplianceConnectionsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) NetworkVirtualApplianceConnectionsListComplete(ctx context.Context, id NetworkVirtualApplianceId) (NetworkVirtualApplianceConnectionsListCompleteResult, error) {
	return c.NetworkVirtualApplianceConnectionsListCompleteMatchingPredicate(ctx, id, NetworkVirtualApplianceConnectionOperationPredicate{})
}

// NetworkVirtualApplianceConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) NetworkVirtualApplianceConnectionsListCompleteMatchingPredicate(ctx context.Context, id NetworkVirtualApplianceId, predicate NetworkVirtualApplianceConnectionOperationPredicate) (result NetworkVirtualApplianceConnectionsListCompleteResult, err error) {
	items := make([]NetworkVirtualApplianceConnection, 0)

	resp, err := c.NetworkVirtualApplianceConnectionsList(ctx, id)
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

	result = NetworkVirtualApplianceConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
