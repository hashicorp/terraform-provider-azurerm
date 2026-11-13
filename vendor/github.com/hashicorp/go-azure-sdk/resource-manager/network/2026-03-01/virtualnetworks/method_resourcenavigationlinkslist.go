package virtualnetworks

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

type ResourceNavigationLinksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceNavigationLink
}

type ResourceNavigationLinksListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceNavigationLink
}

type ResourceNavigationLinksListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ResourceNavigationLinksListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ResourceNavigationLinksList ...
func (c VirtualNetworksClient) ResourceNavigationLinksList(ctx context.Context, id commonids.SubnetId) (result ResourceNavigationLinksListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ResourceNavigationLinksListCustomPager{},
		Path:       fmt.Sprintf("%s/resourceNavigationLinks", id.ID()),
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
		Values *[]ResourceNavigationLink `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ResourceNavigationLinksListComplete retrieves all the results into a single object
func (c VirtualNetworksClient) ResourceNavigationLinksListComplete(ctx context.Context, id commonids.SubnetId) (ResourceNavigationLinksListCompleteResult, error) {
	return c.ResourceNavigationLinksListCompleteMatchingPredicate(ctx, id, ResourceNavigationLinkOperationPredicate{})
}

// ResourceNavigationLinksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworksClient) ResourceNavigationLinksListCompleteMatchingPredicate(ctx context.Context, id commonids.SubnetId, predicate ResourceNavigationLinkOperationPredicate) (result ResourceNavigationLinksListCompleteResult, err error) {
	items := make([]ResourceNavigationLink, 0)

	resp, err := c.ResourceNavigationLinksList(ctx, id)
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

	result = ResourceNavigationLinksListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
