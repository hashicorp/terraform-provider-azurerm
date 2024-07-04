package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayCustomDomainsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GatewayCustomDomainResource
}

type GatewayCustomDomainsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GatewayCustomDomainResource
}

type GatewayCustomDomainsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GatewayCustomDomainsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GatewayCustomDomainsList ...
func (c AppPlatformClient) GatewayCustomDomainsList(ctx context.Context, id GatewayId) (result GatewayCustomDomainsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GatewayCustomDomainsListCustomPager{},
		Path:       fmt.Sprintf("%s/domains", id.ID()),
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
		Values *[]GatewayCustomDomainResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GatewayCustomDomainsListComplete retrieves all the results into a single object
func (c AppPlatformClient) GatewayCustomDomainsListComplete(ctx context.Context, id GatewayId) (GatewayCustomDomainsListCompleteResult, error) {
	return c.GatewayCustomDomainsListCompleteMatchingPredicate(ctx, id, GatewayCustomDomainResourceOperationPredicate{})
}

// GatewayCustomDomainsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) GatewayCustomDomainsListCompleteMatchingPredicate(ctx context.Context, id GatewayId, predicate GatewayCustomDomainResourceOperationPredicate) (result GatewayCustomDomainsListCompleteResult, err error) {
	items := make([]GatewayCustomDomainResource, 0)

	resp, err := c.GatewayCustomDomainsList(ctx, id)
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

	result = GatewayCustomDomainsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
