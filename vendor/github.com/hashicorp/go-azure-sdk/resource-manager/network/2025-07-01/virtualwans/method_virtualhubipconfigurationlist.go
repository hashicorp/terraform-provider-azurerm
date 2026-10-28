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

type VirtualHubIPConfigurationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HubIPConfiguration
}

type VirtualHubIPConfigurationListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HubIPConfiguration
}

type VirtualHubIPConfigurationListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VirtualHubIPConfigurationListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VirtualHubIPConfigurationList ...
func (c VirtualWANsClient) VirtualHubIPConfigurationList(ctx context.Context, id VirtualHubId) (result VirtualHubIPConfigurationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VirtualHubIPConfigurationListCustomPager{},
		Path:       fmt.Sprintf("%s/ipConfigurations", id.ID()),
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
		Values *[]HubIPConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualHubIPConfigurationListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualHubIPConfigurationListComplete(ctx context.Context, id VirtualHubId) (VirtualHubIPConfigurationListCompleteResult, error) {
	return c.VirtualHubIPConfigurationListCompleteMatchingPredicate(ctx, id, HubIPConfigurationOperationPredicate{})
}

// VirtualHubIPConfigurationListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualHubIPConfigurationListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate HubIPConfigurationOperationPredicate) (result VirtualHubIPConfigurationListCompleteResult, err error) {
	items := make([]HubIPConfiguration, 0)

	resp, err := c.VirtualHubIPConfigurationList(ctx, id)
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

	result = VirtualHubIPConfigurationListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
