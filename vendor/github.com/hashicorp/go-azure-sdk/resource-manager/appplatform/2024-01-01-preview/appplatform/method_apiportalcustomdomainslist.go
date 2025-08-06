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

type ApiPortalCustomDomainsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiPortalCustomDomainResource
}

type ApiPortalCustomDomainsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiPortalCustomDomainResource
}

type ApiPortalCustomDomainsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApiPortalCustomDomainsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApiPortalCustomDomainsList ...
func (c AppPlatformClient) ApiPortalCustomDomainsList(ctx context.Context, id ApiPortalId) (result ApiPortalCustomDomainsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ApiPortalCustomDomainsListCustomPager{},
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
		Values *[]ApiPortalCustomDomainResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ApiPortalCustomDomainsListComplete retrieves all the results into a single object
func (c AppPlatformClient) ApiPortalCustomDomainsListComplete(ctx context.Context, id ApiPortalId) (ApiPortalCustomDomainsListCompleteResult, error) {
	return c.ApiPortalCustomDomainsListCompleteMatchingPredicate(ctx, id, ApiPortalCustomDomainResourceOperationPredicate{})
}

// ApiPortalCustomDomainsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) ApiPortalCustomDomainsListCompleteMatchingPredicate(ctx context.Context, id ApiPortalId, predicate ApiPortalCustomDomainResourceOperationPredicate) (result ApiPortalCustomDomainsListCompleteResult, err error) {
	items := make([]ApiPortalCustomDomainResource, 0)

	resp, err := c.ApiPortalCustomDomainsList(ctx, id)
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

	result = ApiPortalCustomDomainsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
