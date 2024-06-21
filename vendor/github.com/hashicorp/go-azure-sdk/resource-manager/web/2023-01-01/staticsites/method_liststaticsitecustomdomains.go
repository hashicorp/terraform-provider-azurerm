package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListStaticSiteCustomDomainsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteCustomDomainOverviewARMResource
}

type ListStaticSiteCustomDomainsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteCustomDomainOverviewARMResource
}

// ListStaticSiteCustomDomains ...
func (c StaticSitesClient) ListStaticSiteCustomDomains(ctx context.Context, id StaticSiteId) (result ListStaticSiteCustomDomainsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/customDomains", id.ID()),
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
		Values *[]StaticSiteCustomDomainOverviewARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListStaticSiteCustomDomainsComplete retrieves all the results into a single object
func (c StaticSitesClient) ListStaticSiteCustomDomainsComplete(ctx context.Context, id StaticSiteId) (ListStaticSiteCustomDomainsCompleteResult, error) {
	return c.ListStaticSiteCustomDomainsCompleteMatchingPredicate(ctx, id, StaticSiteCustomDomainOverviewARMResourceOperationPredicate{})
}

// ListStaticSiteCustomDomainsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) ListStaticSiteCustomDomainsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteCustomDomainOverviewARMResourceOperationPredicate) (result ListStaticSiteCustomDomainsCompleteResult, err error) {
	items := make([]StaticSiteCustomDomainOverviewARMResource, 0)

	resp, err := c.ListStaticSiteCustomDomains(ctx, id)
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

	result = ListStaticSiteCustomDomainsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
