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

type ServiceAssociationLinksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceAssociationLink
}

type ServiceAssociationLinksListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServiceAssociationLink
}

type ServiceAssociationLinksListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServiceAssociationLinksListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServiceAssociationLinksList ...
func (c VirtualNetworksClient) ServiceAssociationLinksList(ctx context.Context, id commonids.SubnetId) (result ServiceAssociationLinksListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServiceAssociationLinksListCustomPager{},
		Path:       fmt.Sprintf("%s/serviceAssociationLinks", id.ID()),
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
		Values *[]ServiceAssociationLink `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServiceAssociationLinksListComplete retrieves all the results into a single object
func (c VirtualNetworksClient) ServiceAssociationLinksListComplete(ctx context.Context, id commonids.SubnetId) (ServiceAssociationLinksListCompleteResult, error) {
	return c.ServiceAssociationLinksListCompleteMatchingPredicate(ctx, id, ServiceAssociationLinkOperationPredicate{})
}

// ServiceAssociationLinksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworksClient) ServiceAssociationLinksListCompleteMatchingPredicate(ctx context.Context, id commonids.SubnetId, predicate ServiceAssociationLinkOperationPredicate) (result ServiceAssociationLinksListCompleteResult, err error) {
	items := make([]ServiceAssociationLink, 0)

	resp, err := c.ServiceAssociationLinksList(ctx, id)
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

	result = ServiceAssociationLinksListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
