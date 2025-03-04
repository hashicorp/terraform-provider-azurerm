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

type GetUserProvidedFunctionAppsForStaticSiteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteUserProvidedFunctionAppARMResource
}

type GetUserProvidedFunctionAppsForStaticSiteCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteUserProvidedFunctionAppARMResource
}

type GetUserProvidedFunctionAppsForStaticSiteCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetUserProvidedFunctionAppsForStaticSiteCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetUserProvidedFunctionAppsForStaticSite ...
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSite(ctx context.Context, id StaticSiteId) (result GetUserProvidedFunctionAppsForStaticSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetUserProvidedFunctionAppsForStaticSiteCustomPager{},
		Path:       fmt.Sprintf("%s/userProvidedFunctionApps", id.ID()),
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
		Values *[]StaticSiteUserProvidedFunctionAppARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetUserProvidedFunctionAppsForStaticSiteComplete retrieves all the results into a single object
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSiteComplete(ctx context.Context, id StaticSiteId) (GetUserProvidedFunctionAppsForStaticSiteCompleteResult, error) {
	return c.GetUserProvidedFunctionAppsForStaticSiteCompleteMatchingPredicate(ctx, id, StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate{})
}

// GetUserProvidedFunctionAppsForStaticSiteCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetUserProvidedFunctionAppsForStaticSiteCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteUserProvidedFunctionAppARMResourceOperationPredicate) (result GetUserProvidedFunctionAppsForStaticSiteCompleteResult, err error) {
	items := make([]StaticSiteUserProvidedFunctionAppARMResource, 0)

	resp, err := c.GetUserProvidedFunctionAppsForStaticSite(ctx, id)
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

	result = GetUserProvidedFunctionAppsForStaticSiteCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
